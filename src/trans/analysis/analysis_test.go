package analysis_test

import (
	"fmt"
	"testing"
	"trans/analysis"
	"trans/filetool"
)

func Test_example(t *testing.T) {
	text, err := filetool.GetInstance().ReadAll("test.lua")
	if err != nil {
		t.Fatal(err)
	}
	ana := analysis.New()
	fanalysis, _, err := ana.GetTool("test.lua")
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

func Benchmark_example(b *testing.B) {
	text, err := filetool.GetInstance().ReadAll("test.lua")
	if err != nil {
		b.Fatal("can not read file")
	}
	ana := analysis.New()
	fanalysis, _, err := ana.GetTool("test.lua")
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
