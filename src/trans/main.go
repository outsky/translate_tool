package main

import (
	"fmt"
	"strings"
	"trans/analysis"
	"trans/cmd"
	"trans/filetool"
	"trans/log"
)

const (
	const_config_file string = "config.ini"
	const_ignore_file string = "ignore.conf"
)

func initConfig() {
	ft := filetool.GetInstance()
	bv, err := ft.ReadFileLine(const_config_file)
	if err != nil {
		bv = [][]byte{
			[]byte(";通过文件扩展名配置提取规则"),
			[]byte(";支持‘lua_rules’，‘prefab_rules’，‘table_rules’, 'ini_rules'"),
			[]byte("[rules]"),
			[]byte(".lua=lua_rules"),
			[]byte(".prefab=prefab_rules"),
			[]byte(".tab=table_rules"),
			[]byte(";根据文件扩展名设置文件读取编码"),
			[]byte(";支持utf8，gbk，hz-gb2312，gb18030，big5"),
			[]byte("[encode]"),
			[]byte(".lua=utf8"),
			[]byte(".prefab=utf8"),
			[]byte(".tab=gbk"),
		}
		err = ft.SaveFileLine(const_config_file, bv)
		if err != nil {
			log.Error("fc", err)
		}
	}
	anal := analysis.GetInstance()
	var nType int
	for _, v := range bv {
		if v[0] == byte(';') {
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
		default:
			switch nType {
			case 1:
				kv := strings.Split(s, "=")
				if len(kv) < 2 {
					panic(fmt.Sprintf("config error: %s", s))
				}
				anal.SetRulesMap(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
			case 2:
				kv := strings.Split(s, "=")
				if len(kv) < 2 {
					panic(fmt.Sprintf("config error: %s", s))
				}
				ft.SetEncoding(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
			}
		}
	}
}

func initFilter() {
	ft := filetool.GetInstance()
	bv, err := ft.ReadFileLine(const_ignore_file)
	if err != nil {
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
			log.Error("fc", err)
		}
	}
	anal := analysis.GetInstance()
	for _, v := range bv {
		if v[0] == byte(';') {
			continue
		}
		s := string(v)
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			anal.SetFilterMap(s)
		}
	}
}

func main() {
	defer log.CloseLog()

	initConfig()
	initFilter()

	//init cobra
	cmd.Execute()
}
