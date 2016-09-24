package common_test

import (
	"fmt"
	"testing"
	"trans/analysis/common"
	"trans/filetool"
)

func Test_Example(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding("html", "gbk")
	text, err := ft.ReadAll("../../test/cn/test.html")
	if err != nil {
		t.Fatal(err)
	}
	entry, start, end, err := common.New().GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d [%d:%d] %s\n", i+1, start[i], end[i], entry[i])
	}
}
