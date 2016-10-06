package log

import (
	_log "log"
	"os"
	"strings"
	"sync"
)

const (
	flagFile = 1 << iota
	flagConsole
)

const (
	preInfo  string = "   "
	preError string = "[x]"
)

type log struct {
	fhandle *os.File
	file    *_log.Logger
	console *_log.Logger
}

var errlog error
var logpath string

var instance *log
var once sync.Once

func getinstance() *log {
	once.Do(func() {
		instance = &log{}
		if len(logpath) <= 0 {
			logpath = "log.txt"
		}
		instance.fhandle, errlog = os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if errlog != nil {
			panic(errlog)
		}
		instance.file = _log.New(instance.fhandle, "", _log.Lshortfile|_log.Ltime)
		instance.console = _log.New(os.Stdout, "", 0)
	})
	return instance
}

func writeLog(flag int, prefix string, msg string) {
	l := getinstance()
	if l.file != nil && flag&flagFile != 0 {
		l.file.SetPrefix(prefix)
		l.file.Println(msg)
	}
	if l.console != nil && flag&flagConsole != 0 {
		l.console.SetPrefix(prefix)
		l.console.Println(msg)
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

func (l *log) Close() {
	if l.fhandle != nil {
		l.fhandle.Close()
	}
}

func SetLogPath(path string) {
	logpath = strings.TrimSpace(path)
}

func Info(flag string, msg string) {
	f := parseFlag(flag)
	writeLog(f, preInfo, msg)
}

func Error(flag string, msg string) {
	f := parseFlag(flag)
	writeLog(f, preError, msg)
}
