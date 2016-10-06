package ini_test

import (
	"fmt"
	"testing"

	"trans/analysis/ini"
	"trans/filetool"
)

func Test_Example(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding(".tab", "gbk")
	text, err := ft.ReadAll("../../test/cn/test.ini")
	if err != nil {
		t.Fatal(err)
	}
	entry, start, end, err := ini.New().GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d [%d:%d] %s\n", i+1, start[i], end[i], entry[i])
	}
}
