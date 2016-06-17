package main

import (
	"log"
	"os"
	"testing"
)

func Test_Example(t *testing.T) {
	flog, err := os.Create(const_log_file)
	if err != nil {
		log.Fatalln(err)
	}
	defer flog.Close()
	logFile = log.New(flog, "", log.LstdFlags)
	logPrint = log.New(os.Stdout, "", log.LstdFlags)
	rulesMap = map[string]string{"lua": "lua_rules", "prefab": "prefab_rules", "tab": "table_rules", "txt": "table_rules"}
	encodeMap = map[string]string{"lua": "utf8", "prefab": "utf8", "tab": "gbk", "txt": "gbk"}
	filterMap = map[string]bool{"filter": true, "filter.lua": true}
	filterExtension = []string{"lua", "prefab", "tab"}
	GetString("test/cn")
	Update("chinese.txt", "test/trans.txt")
	Translate("test/cn", "test/en", 1)
}
