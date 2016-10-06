package log

import (
	_log "log"
	"os"
	"sync"
)

const (
	flagFile = 1 << iota
	flagConsole
)

const (
	preInfo  string = "[INFO]  "
	preError string = "[ERROR] "
)

type log struct {
	fhandle *os.File
	file    *_log.Logger
	console *_log.Logger
}

var errlog error

const (
	const_log_file string = "log.txt"
)

var instance *log
var once sync.Once

func getinstance() *log {
	once.Do(func() {
		instance = &log{}
		instance.fhandle, errlog = os.Create(const_log_file)
		if errlog != nil {
			panic(errlog)
		}
		instance.file = _log.New(instance.fhandle, "", _log.LstdFlags)
		instance.console = _log.New(os.Stdout, "", _log.LstdFlags)
	})
	return instance
}

func (l *log) Close() {
	if l.fhandle != nil {
		l.fhandle.Close()
	}
}

func writeLog(flag int, prefix string, v ...interface{}) {
	l := getinstance()
	if l.file != nil && flag&flagFile != 0 {
		l.file.SetPrefix(prefix)
		l.file.Println(v...)
	}
	if l.console != nil && flag&flagConsole != 0 {
		l.console.SetPrefix(prefix)
		l.console.Println(v...)
	}
}

// flag:
//	f - file
//	c - console
func parseFlag(flag string) int {
	ret := 0
	for _, c := range flag {
		if c == 'f' {
			ret |= flagFile
		} else if c == 'c' {
			ret |= flagConsole
		}
	}
	return ret
}

func CloseLog() {
	getinstance().Close()
}

func Info(flag string, v ...interface{}) {
	f := parseFlag(flag)
	writeLog(f, preInfo, v)
}

func Error(flag string, v ...interface{}) {
	f := parseFlag(flag)
	writeLog(f, preError, v)
}
