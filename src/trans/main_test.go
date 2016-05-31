package main

import (
	"fmt"
	"testing"
	"trans/dic"
)

func Test_GetString(t *testing.T) {
	GetString("test")
	GetString("f:/bqp/bqp/client1")
}

func Test_Update(t *testing.T) {
	Update("test/cn.txt", "test/en.txt")
	db := dic.New("dictionary.db")
	defer db.Close()
	ret, err := db.Query([]byte("家族排名："))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", ret)
	Update("test/cn.txt", "test/test.lua")
}

func Test_Translate(t *testing.T) {
	Translate("test/cn", "test/en", 1)
}

func Benchmark_GetString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetString("test")
	}
}

func Benchmark_Update(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Update("test/cn.txt", "test/en.txt")
	}
}
