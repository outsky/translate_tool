package analysis

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"trans/analysis/common"
	"trans/analysis/ini"
	"trans/analysis/lua"
	"trans/analysis/prefab"
	"trans/analysis/tabfile"
	"trans/dic"
	"trans/filetool"
	"trans/gpool"
	"trans/log"
)

type delegate interface {
	GetString(text []byte) ([][]byte, []int, []int, error)
	Pretreat(trans []byte) []byte
}

type analysis struct {
	rulesMap  map[string]string
	filterMap map[string]bool
}

const (
	const_rule_common    = "common_rules"
	const_rule_lua       = "lua_rules"
	const_rule_prefab    = "prefab_rules"
	const_rule_tablefile = "table_rules"
	const_rule_ini       = "ini_rules"
)

var instance *analysis
var once sync.Once

func GetInstance() *analysis {
	once.Do(func() {
		instance = &analysis{
			rulesMap:  make(map[string]string),
			filterMap: make(map[string]bool),
		}
	})
	return instance
}

func (a *analysis) SetRulesMap(k, v string) {
	a.rulesMap[path.Ext(k)] = v
}

func (a *analysis) SetFilterMap(key string) {
	a.filterMap[key] = true
}

func (a *analysis) getPool(file string) (delegate, error) {
	file_ex := path.Ext(file)
	rule, ok := a.rulesMap[file_ex]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no rule: %s", file))
	}
	switch rule {
	case const_rule_common:
		return common.New(file), nil
	case const_rule_lua:
		return lua.New(file), nil
	case const_rule_prefab:
		return prefab.New(file), nil
	case const_rule_tablefile:
		return tabfile.New(file), nil
	case const_rule_ini:
		return ini.New(file), nil
	default:
		return nil, errors.New(fmt.Sprintf("rule not defined: %s(%s)", file, rule))
	}
}

func (a *analysis) shouldIgnore(name string) bool {
	ext := path.Ext(name)
	if _, ok := a.filterMap[ext]; ok {
		return true
	}
	namev := strings.Split(name, "/")
	for _, filename := range namev {
		if _, ok := a.filterMap[filename]; ok {
			return true
		}
	}
	return false
}

func (a *analysis) GetString(dbname, update, root string) {
	root = strings.TrimRight(strings.Replace(root, "\\", "/", -1), "/")
	log.Info("fc", fmt.Sprintf("extract chinese from %s", root))
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(root)
	if err != nil {
		log.Error("fc", err.Error())
		return
	}
	newcount := 0
	db := dic.NewDic(dbname)
	notrans := dic.NewUpt(update)
	for i := 0; i < len(fmap); i++ {
		if a.shouldIgnore(fmap[i]) {
			continue
		}
		ins, err := a.getPool(fmap[i])
		if err != nil {
			log.Info("f", err.Error())
			continue
		}
		context, err := ft.ReadAll(fmap[i])
		if err != nil {
			log.Info("fc", err.Error())
			continue
		}
		entry, _, _, err := ins.GetString(context)
		if err != nil {
			log.Error("fc", err.Error())
		}
		relativepath := strings.Split(fmap[i], root)[1]
		for _, v := range entry {
			if _, ok := db.Query(v); !ok {
				if notrans.Append(v, []byte(""), relativepath) {
					newcount += 1
				}
			}
		}
	}
	if newcount > 0 {
		notrans.Save()
		log.Info("fc", fmt.Sprintf("generate %s, new line number: %d. finished!", update, newcount))
	} else {
		log.Info("fc", fmt.Sprintf("nothing to do. finished!"))
	}
}

func (a *analysis) Translate(dbname, update, root, output string, queue int, logpath string) {
	log.SetLogPath(logpath)
	root = strings.TrimRight(strings.Replace(root, "\\", "/", -1), "/")
	output = strings.TrimRight(strings.Replace(output, "\\", "/", -1), "/")
	log.Info("fc", fmt.Sprintf("translate %s to %s", root, output))
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(root)
	if err != nil {
		log.Error("fc", err.Error())
		return
	}
	dbdata := dic.NewDic(dbname)
	notrans := dic.NewUpt(update)
	copycount, transcount, newcount, ignorecount := 0, 0, 0, 0
	pool := gpool.New(queue)
	mutex := &sync.Mutex{}
	fwork := func(oldfile, newfile, relativepath string) {
		defer pool.Done()
		var (
			entry   [][]byte
			start   []int
			end     []int
			context [][]byte
			nStart  int
			nSize   int
		)
		bv, err := ft.ReadAll(oldfile)
		if err != nil {
			log.Error("fc", err.Error())
			return
		}
		if a.shouldIgnore(oldfile) {
			ignorecount += 1
			return
		}
		ins, err := a.getPool(oldfile)
		if err != nil {
			log.Info("f", err.Error())
			goto Point
		}
		entry, start, end, err = ins.GetString(bv)
		if err != nil {
			log.Error("fc", err.Error())
			goto Point
		}
		nStart = 0
		nSize = len(bv)
		for i := 0; i < len(entry); i++ {
			context = append(context, bv[nStart:start[i]])
			nStart = end[i]
			if trans, ok := dbdata.Query(entry[i]); ok {
				if len(trans) > 0 {
					context = append(context, ins.Pretreat(trans))
				} else {
					context = append(context, bv[start[i]:end[i]])
					mutex.Lock()
					if notrans.Append(entry[i], []byte(""), relativepath) {
						newcount += 1
					}
					mutex.Unlock()
				}
			} else {
				context = append(context, bv[start[i]:end[i]])
				mutex.Lock()
				if notrans.Append(entry[i], []byte(""), relativepath) {
					newcount += 1
				}
				mutex.Unlock()
			}
		}
		if nStart < nSize {
			context = append(context, bv[nStart:nSize])
		}
	Point:
		if len(context) > 0 {
			oldencoding, err := ft.SetEncoding(newfile, "utf8")
			if err != nil {
				log.Error("fc", err.Error())
			} else {
				if err := ft.WriteAll(newfile, bytes.Join(context, []byte(""))); err != nil {
					log.Error("fc", err.Error())
				} else {
					transcount += 1
				}
				ft.SetEncoding(newfile, oldencoding)
			}
		} else {
			if err := ft.WriteAll(newfile, bv); err != nil {
				log.Error("fc", err.Error())
			} else {
				copycount += 1
			}
		}
	}
	for i := 0; i < len(fmap); i++ {
		pool.Add(1)
		relativepath := strings.Split(fmap[i], root)[1]
		fpath := strings.Replace(fmap[i], root, output, 1)
		go fwork(fmap[i], fpath, relativepath)
	}
	pool.Wait()
	if newcount > 0 {
		notrans.Save()
		log.Info("fc", fmt.Sprintf("generate %s, new line number: %d.", update, newcount))
	}
	log.Info("fc", fmt.Sprintf("translate file %d, copy file %d, ignore file %d, total %d/%d.\n\n",
		transcount, copycount, ignorecount, transcount+copycount+ignorecount, len(fmap)))
	return
}

func (a *analysis) Update(dbname, update string) {
	log.Info("fc", fmt.Sprintf("update %s to %s", update, dbname))
	dbdata := dic.NewDic(dbname)
	newdata := dic.NewDic(update)
	text, trans := newdata.GetLine()
	count := len(text)
	if count > 0 {
		for i := 0; i < count; i++ {
			dbdata.Append(text[i], trans[i])
		}
		dbdata.Save()
		log.Info("fc", fmt.Sprintf("update line number: %d. finished!", count))
	} else {
		log.Info("fc", "nothing to do. finished!")
	}
}
