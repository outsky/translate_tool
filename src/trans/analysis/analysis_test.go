package analysis_test

import (
	"fmt"
	"testing"
	"trans/analysis"
	"trans/filetool"
)

func Test_lua(t *testing.T) {
	ana := analysis.GetInstance()
	ana.SetFilterFileEx([]string{"lua", "tab", "prefab"})
	ana.SetRule(map[string]string{"lua": "lua_rules", "prefab": "prefab_rules", "tab": "table_rules"})
	fanalysis, _, err := ana.GetRule("../test/test.lua")
	if err != nil {
		t.Fatal(err)
	}
	filetool.GetInstance().SetEncoding("lua", "utf8")
	text, err := filetool.GetInstance().ReadAll("../test/test.lua")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := fanalysis(text)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range entry {
		fmt.Printf("%s\n", v)
	}
}

func Test_prefab(t *testing.T) {
	ana := analysis.GetInstance()
	fanalysis, _, err := ana.GetRule("../test/cn/Boss.prefab")
	if err != nil {
		t.Fatal(err)
	}
	filetool.GetInstance().SetEncoding("prefab", "utf8")
	text, err := filetool.GetInstance().ReadAll("../test/cn/Boss.prefab")
	if err != nil {
		t.Fatal(err)
	}
	entry, err1 := fanalysis(text)
	if err1 != nil {
		t.Fatal(err)
	}
	for k, v := range entry {
		fmt.Printf("%d %s\n", k, v)
	}
}

func Test_tab(t *testing.T) {
	ana := analysis.GetInstance()
	fanalysis, _, err := ana.GetRule("../test/cn/ScriptItem.tab")
	if err != nil {
		t.Fatal(err)
	}
	filetool.GetInstance().SetEncoding("tab", "gbk")
	text, err := filetool.GetInstance().ReadAll("../test/cn/ScriptItem.tab")
	if err != nil {
		t.Fatal(err)
	}
	entry, err1 := fanalysis(text)
	if err1 != nil {
		t.Fatal(err)
	}
	for k, v := range entry {
		fmt.Printf("%d %s\n", k, v)
	}
}

func Benchmark_lua(b *testing.B) {
	ana := analysis.GetInstance()
	fanalysis, _, err := ana.GetRule("../test/test.lua")
	if err != nil {
		b.Fatal(err)
	}
	filetool.GetInstance().SetEncoding("lua", "utf8")
	text, err := filetool.GetInstance().ReadAll("../test/test.lua")
	if err != nil {
		b.Fatal("can not read file")
	}
	for i := 0; i < b.N; i++ {
		_, err = fanalysis(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}
