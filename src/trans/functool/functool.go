package functool

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"trans/analysis"
	"trans/dic"
	"trans/filetool"
	"trans/gpool"
)

const (
	const_ignore_file  string = "ignore.conf"
	const_chinese_file string = "chinese.txt"
	const_dic_file     string = "dictionary.db"
	const_log_file     string = "log.txt"
)

var filterMap map[string]string
var logFile *log.Logger
var logPrint *log.Logger

const (
	log_file = 1 << iota
	log_print
)

func writeLog(flag int, v ...interface{}) {
	if flag&log_file != 0 {
		logFile.Println(v...)
	}
	if flag&log_print != 0 {
		logPrint.Println(v...)
	}
}

func init() {
	// create logger
	flog, err := os.Create(const_log_file)
	if err != nil {
		log.Fatalln(err)
	}
	logFile = log.New(flog, "[trans]", log.LstdFlags)
	logPrint = log.New(os.Stdout, "[trans]", log.LstdFlags)
	// init or read ignore config
	filterMap = make(map[string]string)
	ft := filetool.GetInstance()
	bv, err := ft.ReadFileLine(const_ignore_file)
	if err != nil {
		writeLog(log_file, err)
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
			writeLog(log_file|log_print, err)
		}
	}
	for _, v := range bv {
		if v[0] == 0x3b {
			continue
		}
		v = bytes.Trim(v, " ")
		sv := string(v)
		filterMap[sv] = sv
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
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(filedir)
	if err != nil {
		writeLog(log_file|log_print, err)
		return
	}
	var entry_total [][]byte
	anal := analysis.New()
	for i := 0; i < len(fmap); i++ {
		if err := filterFile(fmap[i]); err != nil {
			writeLog(log_file, err)
			continue
		}
		fanalysis, _, bDecoder, err := anal.GetRule(fmap[i])
		if err != nil {
			writeLog(log_file, err)
			continue
		}
		context, err := ft.ReadAll(fmap[i], bDecoder)
		if err != nil {
			writeLog(log_file|log_print, err)
			continue
		}
		entry, err := fanalysis(&context)
		if err != nil {
			writeLog(log_file|log_print, err)
		}
		for _, v := range *entry {
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
	if err := ft.SaveFileLine(const_chinese_file, ret); err != nil {
		writeLog(log_file|log_print, err)
		return
	}
	writeLog(log_file|log_print,
		fmt.Sprintf("generate %s, view it. getstring finished!", const_chinese_file))
	return
}

func Update(cnFile, transFile string) {
	cnFile = strings.Replace(cnFile, "\\", "/", -1)
	transFile = strings.Replace(transFile, "\\", "/", -1)
	ft := filetool.GetInstance()
	cnText, err1 := ft.ReadFileLine(cnFile)
	if err1 != nil {
		writeLog(log_file|log_print, err1)
		return
	}
	transText, err2 := ft.ReadFileLine(transFile)
	if err2 != nil {
		writeLog(log_file|log_print, err2)
		return
	}
	cnLen := len(cnText)
	transLen := len(transText)
	if cnLen != transLen {
		writeLog(log_file|log_print, fmt.Sprintf("line number is not equal: %s:%d %s:%d",
			cnFile, cnLen, transFile, transLen))
		return
	}
	db := dic.New(const_dic_file)
	defer db.Close()
	var count int = 0
	for i := 0; i < cnLen; i++ {
		if err := db.Insert(cnText[i], transText[i]); err != nil {
			writeLog(log_file|log_print,
				fmt.Sprintf("insert to db failed: %s:%s", cnText[i], transText[i]))
		} else {
			count++
		}
	}
	if count != cnLen {
		writeLog(log_file|log_print,
			fmt.Sprintf("only insert %d line to dic, total %d.", count, cnLen))
		return
	}
	writeLog(log_file|log_print,
		fmt.Sprintf("update %d/%d line number to dic. update finished!", count, cnLen))
	return
}

func Translate(src, des string, queue int) {
	src = strings.Replace(src, "\\", "/", -1)
	des = strings.Replace(des, "\\", "/", -1)
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(src)
	if err != nil {
		writeLog(log_file|log_print, err)
		return
	}
	writeLog(log_file|log_print, fmt.Sprintf("translate files total: %d", len(fmap)))
	db := dic.New(const_dic_file)
	defer db.Close()
	var notrans [][]byte
	pool := gpool.New(queue)
	mutex := &sync.Mutex{}
	f := func(oldfile, newfile string) {
		defer pool.Done()
		var entry *[][]byte
		anal := analysis.New()
		fanalysis, ftranslate, bDecoder, err := anal.GetRule(oldfile)
		bv, e := ft.ReadAll(oldfile, bDecoder)
		if e != nil {
			writeLog(log_file|log_print, e)
			return
		}
		if err != nil {
			writeLog(log_file, err)
			goto Point
		}
		if err = filterFile(oldfile); err != nil {
			writeLog(log_file, err)
			goto Point
		}
		entry, err = fanalysis(&bv)
		if err != nil {
			writeLog(log_file|log_print, err)
			goto Point
		}
		for _, v := range *entry {
			trans, err := db.Query(v)
			if err != nil {
				mutex.Lock()
				notrans = append(notrans, v)
				mutex.Unlock()
				continue
			}
			if err := ftranslate(&bv, v, trans); err != nil {
				writeLog(log_file|log_print, err)
			}
		}
	Point:
		ft.WriteAll(newfile, bv, bDecoder)
	}
	for i := 0; i < len(fmap); i++ {
		fpath := strings.Replace(fmap[i], src, des, 1)
		pool.Add(1)
		go f(fmap[i], fpath)
	}
	pool.Wait()
	if len(notrans) > 0 {
		if err := ft.SaveFileLine(const_chinese_file, notrans); err != nil {
			writeLog(log_file|log_print, err)
			return
		}
		writeLog(log_file|log_print,
			fmt.Sprintf("update %s, line number: %d.", const_chinese_file, len(notrans)))
	}
	writeLog(log_file|log_print, "translate finished!")
	return
}
