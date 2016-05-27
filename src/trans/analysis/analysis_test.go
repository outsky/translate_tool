package analysis_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"trans/analysis"
	"trans/filetool"

	"github.com/rainycape/unidecode"
)

func Test_lua(t *testing.T) {
	text, err := filetool.GetInstance().ReadAll("test.lua")
	if err != nil {
		t.Fatal(err)
	}
	ana := analysis.New()
	fanalysis, _, err := ana.GetRule("test.lua")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := fanalysis(&text)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range *entry {
		fmt.Printf("%s\n", v)
	}
}

func Test_prefab(t *testing.T) {
	text, err := filetool.GetInstance().ReadAll("Boss.prefab")
	if err != nil {
		t.Fatal(err)
	}
	ana := analysis.New()
	fanalysis, _, err := ana.GetRule("Boss.prefab")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := fanalysis(&text)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range *entry {
		fmt.Printf("%d %s %s\n", k, v, unidecode.Unidecode(string(v)))
	}

	s := []byte("abc\r测试!")
	quoted := strconv.QuoteToASCII(string(s))
	unquoted := quoted[1 : len(quoted)-1]
	fmt.Println(quoted, unquoted)

	unicode := "\u5BB6\u65CF"
	v := strings.Split(unicode, "\\u")
	name := fmt.Sprintf("%v", v)
	ouye := name[1 : len(name)-1]
	fmt.Println(ouye)
}

func Benchmark_lua(b *testing.B) {
	text, err := filetool.GetInstance().ReadAll("test.lua")
	if err != nil {
		b.Fatal("can not read file")
	}
	ana := analysis.New()
	fanalysis, _, err := ana.GetRule("test.lua")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = fanalysis(&text)
		if err != nil {
			b.Fatal(err)
		}
	}
}
