package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"trans/analysis"
	"trans/dic"
	"trans/filetool"
	"trans/gpool"

	//	"github.com/pkg/profile"
)

const (
	const_config_file  string = "config.ini"
	const_ignore_file  string = "ignore.conf"
	const_chinese_file string = "chinese.txt"
	const_dic_file     string = "dictionary.db"
	const_log_file     string = "log.txt"
)

var rulesMap map[string]string
var encodeMap map[string]string
var filterMap map[string]bool
var filterExtension []string

var logFile *log.Logger
var logPrint *log.Logger

const (
	log_file = 1 << iota
	log_print
)

const (
	log_info  string = "[INFO]  "
	log_error string = "[ERROR] "
)

func writeLog(flag int, level string, v ...interface{}) {
	if flag&log_file != 0 {
		logFile.SetPrefix(level)
		logFile.Println(v...)
	}
	if flag&log_print != 0 {
		logPrint.SetPrefix(level)
		logPrint.Println(v...)
	}
}

func filterFile(name string) error {
	namev := strings.Split(name, "/")
	for _, filename := range namev {
		if _, ok := filterMap[filename]; ok {
			return errors.New(fmt.Sprintf("[ingnore file] %s", name))
		}
	}
	return nil
}

func GetString(filedir string) {
	filedir = strings.Replace(filedir, "\\", "/", -1)
	writeLog(log_file|log_print, log_info, fmt.Sprintf("extract chinese from %s", filedir))
	ft := filetool.GetInstance()
	for k, v := range encodeMap {
		if err := ft.SetEncoding(k, v); err != nil {
			writeLog(log_file|log_print, log_error, err)
		}
	}
	fmap, err := ft.GetFilesMap(filedir)
	if err != nil {
		writeLog(log_file|log_print, log_error, err)
		return
	}
	var entry_total [][]byte
	anal := analysis.GetInstance()
	anal.SetRule(rulesMap)
	anal.SetFilterFileEx(filterExtension)
	for i := 0; i < len(fmap); i++ {
		if err := filterFile(fmap[i]); err != nil {
			writeLog(log_file|log_print, log_info, err)
			continue
		}
		fanalysis, _, err := anal.GetRule(fmap[i])
		if err != nil {
			writeLog(log_file|log_print, log_info, err)
			continue
		}
		context, err := ft.ReadAll(fmap[i])
		if err != nil {
			writeLog(log_file|log_print, log_info, err)
			continue
		}
		entry, err := fanalysis(context)
		if err != nil {
			writeLog(log_file|log_print, log_error, err)
		}
		for _, v := range entry {
			bIsExsit := false
			for _, m := range entry_total {
				if bytes.Compare(v, m) == 0 {
					bIsExsit = true
				}
			}
			if !bIsExsit {
				entry_total = append(entry_total, v)
			}
		}
	}
	db := dic.New(const_dic_file)
	defer db.Close()
	var ret [][]byte
	for i := 0; i < len(entry_total); i++ {
		if _, err := db.Query(entry_total[i]); err != nil {
			ret = append(ret, entry_total[i])
		}
	}
	ft.SetEncoding(const_chinese_file, "utf8")
	if err := ft.SaveFileLine(const_chinese_file, ret); err != nil {
		writeLog(log_file|log_print, log_error, err)
		return
	}
	writeLog(log_file|log_print, log_info,
		fmt.Sprintf("generate %s, line number: %d. getstring finished!", const_chinese_file, len(ret)))
	return
}

func Update(cnFile, transFile string) {
	cnFile = strings.Replace(cnFile, "\\", "/", -1)
	transFile = strings.Replace(transFile, "\\", "/", -1)
	writeLog(log_file|log_print, log_info, fmt.Sprintf("update dictionary from %s to %s", cnFile, transFile))
	ft := filetool.GetInstance()
	cnText, err1 := ft.ReadFileLine(cnFile)
	if err1 != nil {
		writeLog(log_file|log_print, log_error, err1)
		return
	}
	transText, err2 := ft.ReadFileLine(transFile)
	if err2 != nil {
		writeLog(log_file|log_print, log_error, err2)
		return
	}
	cnLen := len(cnText)
	transLen := len(transText)
	if cnLen != transLen {
		writeLog(log_file|log_print, log_error,
			fmt.Sprintf("line number is not equal: %s:%d %s:%d", cnFile, cnLen, transFile, transLen))
		return
	}
	db := dic.New(const_dic_file)
	defer db.Close()
	for i := 0; i < cnLen; i++ {
		if err := db.Insert(cnText[i], transText[i]); err != nil {
			writeLog(log_file|log_print, log_error,
				fmt.Sprintf("insert to db failed at %d line: %s:%s", i, cnText[i], transText[i]))
		}
	}
	writeLog(log_file|log_print, log_info,
		fmt.Sprintf("update %d line number to %s. update finished!", cnLen, const_dic_file))
	return
}

func Translate(src, des string, queue int) {
	src = strings.TrimRight(strings.Replace(src, "\\", "/", -1), "/")
	des = strings.TrimRight(strings.Replace(des, "\\", "/", -1), "/")
	writeLog(log_file|log_print, log_info, fmt.Sprintf("translate %s to %s", src, des))
	ft := filetool.GetInstance()
	for k, v := range encodeMap {
		if err := ft.SetEncoding(k, v); err != nil {
			writeLog(log_file|log_print, log_error, err)
		}
	}
	fmap, err := ft.GetFilesMap(src)
	if err != nil {
		writeLog(log_file|log_print, log_error, err)
		return
	}
	db := dic.New(const_dic_file)
	defer db.Close()
	var notrans [][]byte
	tatal, transcount := 0, 0
	pool := gpool.New(queue)
	mutex := &sync.Mutex{}
	anal := analysis.GetInstance()
	anal.SetRule(rulesMap)
	anal.SetFilterFileEx(filterExtension)
	f := func(oldfile, newfile string) {
		defer pool.Done()
		var entry [][]byte
		bv, err := ft.ReadAll(oldfile)
		if err != nil {
			writeLog(log_file|log_print, log_error, err)
			return
		}
		fanalysis, ftranslate, err := anal.GetRule(oldfile)
		if err != nil {
			writeLog(log_file|log_print, log_info, err)
			goto Point
		}
		if err = filterFile(oldfile); err != nil {
			writeLog(log_file|log_print, log_info, err)
			goto Point
		}
		entry, err = fanalysis(bv)
		if err != nil {
			writeLog(log_file|log_print, log_error, err)
			goto Point
		}
		for _, v := range entry {
			trans, err := db.Query(v)
			if err != nil {
				bIsExsit := false
				for _, m := range notrans {
					if bytes.Compare(v, m) == 0 {
						bIsExsit = true
					}
				}
				if !bIsExsit {
					mutex.Lock()
					notrans = append(notrans, v)
					mutex.Unlock()
				}
				continue
			}
			if err := ftranslate(&bv, v, trans); err != nil {
				writeLog(log_file|log_print, log_error, err)
			}
		}
		transcount += 1
	Point:
		tatal += 1
		ft.WriteAll(newfile, bv)
	}
	for i := 0; i < len(fmap); i++ {
		fpath := strings.Replace(fmap[i], src, des, 1)
		pool.Add(1)
		go f(fmap[i], fpath)
	}
	pool.Wait()
	if len(notrans) > 0 {
		ft.SetEncoding(const_chinese_file, "utf8")
		if err := ft.SaveFileLine(const_chinese_file, notrans); err != nil {
			writeLog(log_file|log_print, log_error, err)
			return
		}
		writeLog(log_file|log_print, log_info,
			fmt.Sprintf("generate %s, line number: %d.", const_chinese_file, len(notrans)))
	}
	writeLog(log_file|log_print, log_info,
		fmt.Sprintf("translate file %d, copy file %d, finished!", transcount, tatal-transcount))
	return
}

func useage() {
	fmt.Println(
		`trans is a chinese extraction, record and translate tools.

Usage:	trans command [arguments]

The commands are:

	getstring	extract chinese from file or folder.
			e.g. trans getstring path				
	update		update translation to dictionary.
			e.g. trans update chinese.txt translate.txt				
	translate	translate file or folder.
			e.g. trans translate src_path des_path
	
Remark: Supports lua Script, Unity prefab, table file`)
}

func initConfig() {
	ft := filetool.GetInstance()
	bv, err := ft.ReadFileLine(const_config_file)
	if err != nil {
		writeLog(log_file, log_info, err)
		bv = [][]byte{
			[]byte(";通过文件扩展名配置提取规则"),
			[]byte(";支持‘lua_rules’，‘prefab_rules’，‘table_rules’"),
			[]byte("[rules]"),
			[]byte("lua=lua_rules"),
			[]byte("prefab=prefab_rules"),
			[]byte("tab=table_rules"),
			[]byte(";根据文件扩展名设置文件读取编码"),
			[]byte(";支持utf8，gbk，hz-gb2312，gb18030，big5"),
			[]byte("[encode]"),
			[]byte("lua=utf8"),
			[]byte("prefab=utf8"),
			[]byte("tab=gbk"),
			[]byte(";根据文件扩展名判断是否是不翻译需要过滤的路径"),
			[]byte("[filter]"),
			[]byte("extension=lua,prefab,tab"),
		}
		err = ft.SaveFileLine(const_config_file, bv)
		if err != nil {
			writeLog(log_file|log_print, log_error, err)
		}
	}
	rulesMap = make(map[string]string)
	encodeMap = make(map[string]string)
	filterExtension = make([]string, 0)
	var nType int
	for _, v := range bv {
		if v[0] == 0x3b {
			continue
		}
		s := string(v)
		s = strings.TrimSpace(s)
		if len(s) <= 0 {
			continue
		}
		switch s {
		case "[rules]":
			nType = 1
		case "[encode]":
			nType = 2
		case "[filter]":
			nType = 3
		default:
			switch nType {
			case 1:
				kv := strings.Split(s, "=")
				if len(kv) != 2 {
					panic(fmt.Sprintf("config error: %s", s))
				}
				rulesMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			case 2:
				kv := strings.Split(s, "=")
				if len(kv) != 2 {
					panic(fmt.Sprintf("config error: %s", s))
				}
				encodeMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			case 3:
				kv := strings.Split(s, "=")
				if len(kv) != 2 {
					panic(fmt.Sprintf("config error: %s", s))
				}
				exv := strings.Split(kv[1], ",")
				for _, elem := range exv {
					filterExtension = append(filterExtension, strings.TrimSpace(elem))
				}
			}
		}
	}
}

func initFilter() {
	ft := filetool.GetInstance()
	bv, err := ft.ReadFileLine(const_ignore_file)
	if err != nil {
		writeLog(log_file, log_info, err)
		bv = [][]byte{
			[]byte(";这里是忽略的文件，每个文件一行"),
			[]byte(";例如test.lua"),
			[]byte(";自动忽略注释和去空白"),
			[]byte("cvs"),
			[]byte(".git"),
			[]byte(".svn"),
		}
		err = ft.SaveFileLine(const_ignore_file, bv)
		if err != nil {
			writeLog(log_file|log_print, log_error, err)
		}
	}
	filterMap = make(map[string]bool)
	for _, v := range bv {
		if v[0] == 0x3b {
			continue
		}
		s := string(v)
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			filterMap[s] = true
		}
	}
}

func main() {
	//	defer profile.Start(profile.CPUProfile).Stop()
	//	defer profile.Start(profile.MemProfile).Stop()

	// create logger
	flog, err := os.Create(const_log_file)
	if err != nil {
		log.Fatalln(err)
	}
	defer flog.Close()
	logFile = log.New(flog, "", log.LstdFlags)
	logPrint = log.New(os.Stdout, "", log.LstdFlags)

	// init config
	initConfig()

	// init filter file
	initFilter()

	// main
	switch len(os.Args) {
	case 3:
		if strings.EqualFold(os.Args[1], "getstring") {
			GetString(os.Args[2])
		} else {
			useage()
		}
	case 4:
		if strings.EqualFold(os.Args[1], "update") {
			Update(os.Args[2], os.Args[3])
		} else if strings.EqualFold(os.Args[1], "translate") {
			Translate(os.Args[2], os.Args[3], 1)
		} else {
			useage()
		}
	case 5:
		if strings.EqualFold(os.Args[1], "translate") {
			queue, _ := strconv.ParseInt(os.Args[4], 10, 0)
			Translate(os.Args[2], os.Args[3], int(queue))
		} else {
			useage()
		}
	default:
		useage()
	}
}
