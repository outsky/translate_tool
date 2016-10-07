package ini_test

import (
	"fmt"
	"testing"

	"trans/analysis/ini"
	"trans/filetool"
)

func Test_Example(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding(".ini", "gbk")
	text, err := ft.ReadAll("../../test/cn/test2.ini")
	if err != nil {
		t.Fatal(err)
	}
	entry, start, end, err := ini.New("test2").GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d [%d:%d] %s\n", i+1, start[i], end[i], entry[i])
	}
}
