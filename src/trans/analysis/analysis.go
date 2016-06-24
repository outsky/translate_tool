package analysis

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
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
	const_rule_lua       = "lua_rules"
	const_rule_prefab    = "prefab_rules"
	const_rule_tablefile = "table_rules"
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
		return nil, errors.New(fmt.Sprintf("[not extract rule] %s", file))
	}
	switch rule {
	case const_rule_lua:
		return lua.New(), nil
	case const_rule_prefab:
		return prefab.New(), nil
	case const_rule_tablefile:
		return tabfile.New(), nil
	default:
		return nil, errors.New(fmt.Sprintf("[not extract rule] %s", file))
	}
}

func (a *analysis) filter(name string) error {
	namev := strings.Split(name, "/")
	for _, filename := range namev {
		if _, ok := a.filterMap[filename]; ok {
			return errors.New(fmt.Sprintf("[ingnore file] %s", name))
		}
	}
	return nil
}

func (a *analysis) GetString(dbname, update, root string) {
	root = strings.TrimRight(strings.Replace(root, "\\", "/", -1), "/")
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, fmt.Sprintf("extract chinese from %s", root))
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(root)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
		return
	}
	newcount := 0
	db := dic.New(dbname)
	notrans := dic.NewOnly(update)
	for i := 0; i < len(fmap); i++ {
		if err := a.filter(fmap[i]); err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			continue
		}
		ins, err := a.getPool(fmap[i])
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			continue
		}
		context, err := ft.ReadAll(fmap[i])
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			continue
		}
		entry, _, _, err := ins.GetString(context)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
		}
		for _, v := range entry {
			if _, ok := db.Query(v); !ok {
				notrans.Append(v, []byte(""))
				newcount += 1
			}
		}
	}
	if newcount > 0 {
		notrans.Save()
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("generate %s, new line number: %d. finished!", update, newcount))
	} else {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("nothing to do. finished!"))
	}
}

func (a *analysis) Translate(dbname, update, root, output string, queue int) {
	root = strings.TrimRight(strings.Replace(root, "\\", "/", -1), "/")
	output = strings.TrimRight(strings.Replace(output, "\\", "/", -1), "/")
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, fmt.Sprintf("translate %s to %s", root, output))
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(root)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
		return
	}
	dbdata := dic.New(dbname)
	notrans := dic.NewOnly(update)
	tatal, transcount, newcount := 0, 0, 0
	pool := gpool.New(queue)
	mutex := &sync.Mutex{}
	fwork := func(oldfile, newfile string) {
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
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
			return
		}
		ins, err := a.getPool(oldfile)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			goto Point
		}
		if err = a.filter(oldfile); err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			goto Point
		}
		entry, start, end, err = ins.GetString(bv)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
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
					if notrans.Append(entry[i], []byte("")) {
						newcount += 1
					}
					mutex.Unlock()
				}
			} else {
				context = append(context, bv[start[i]:end[i]])
				mutex.Lock()
				if notrans.Append(entry[i], []byte("")) {
					newcount += 1
				}
				newcount += 1
				mutex.Unlock()
			}
		}
		if nStart < nSize {
			context = append(context, bv[nStart:nSize])
		}
		transcount += 1
	Point:
		tatal += 1
		if len(context) > 0 {
			ft.WriteAll(newfile, bytes.Join(context, []byte("")))
		} else {
			ft.WriteAll(newfile, bv)
		}
	}
	for i := 0; i < len(fmap); i++ {
		pool.Add(1)
		fpath := strings.Replace(fmap[i], root, output, 1)
		println(fmap[i], fpath)
		go fwork(fmap[i], fpath)
	}
	pool.Wait()
	if newcount > 0 {
		notrans.Save()
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("generate %s, new line number: %d.", update, newcount))
	}
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
		fmt.Sprintf("translate file %d, copy file %d. finished!", transcount, tatal-transcount))
	return
}

func (a *analysis) Update(dbname, update string) {
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, fmt.Sprintf("update %s to %s", update, dbname))
	dbdata := dic.New(dbname)
	trans := dic.New(update)
	succ := trans.Merge(dbdata)
	if succ > 0 {
		dbdata.Save()
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("update line number: %d. finished!", succ))
	} else {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, "nothing to do. finished!")
	}
}
