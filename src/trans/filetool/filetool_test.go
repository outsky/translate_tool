package filetool_test

import (
	"fmt"
	"testing"
	"trans/filetool"
)

func Test_GetInstance(t *testing.T) {
	ft1 := filetool.GetInstance()
	ft2 := filetool.GetInstance()
	if ft1 != ft2 {
		t.Fatal("GetInstance diffrent value")
	}
}

func Test_GetFileMap(t *testing.T) {
	ft := filetool.GetInstance()
	fm, err := ft.GetFilesMap("../")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(fm); i++ {
		fmt.Println(i, fm[i])
	}
}

func Test_ReadFileLine(t *testing.T) {
	ft := filetool.GetInstance()
	context, err := ft.ReadFileLine("../test/cn.txt")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range context {
		fmt.Printf("%d %s\n", k, v)
	}
}

func Test_SaveFileLine(t *testing.T) {
	ft := filetool.GetInstance()
	context, err := ft.ReadFileLine("../test/cn.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = ft.SaveFileLine("../test/cn2.txt", context)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_ReadAll(t *testing.T) {
	ft := filetool.GetInstance()
	bv, err := ft.ReadAll("../test/test.lua", false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", bv)
}

func Test_WriteAll(t *testing.T) {
	ft := filetool.GetInstance()
	bv, err := ft.ReadAll("../test/test.lua", false)
	if err != nil {
		t.Fatal(err)
	}
	err = ft.WriteAll("../test/test/test/test.lua", bv, false)
	if err != nil {
		t.Fatal(err)
	}
	err = ft.WriteAll("../test/test2.txt", bv, false)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GbkFile(t *testing.T) {
	ft := filetool.GetInstance()
	bv, err := ft.ReadAll("../test/cn/ScriptItem.tab", true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", bv)
	err = ft.WriteAll("../test/ScriptItem2.tab", bv, true)
	if err != nil {
		t.Fatal(err)
	}
}
