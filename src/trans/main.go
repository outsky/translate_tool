package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"trans/functool"

	//	"github.com/pkg/profile"
)

func useage() {
	fmt.Println(
		`trans is a chinese extraction, record and translate tools.

Usage:	trans command [arguments]

The commands are:

	getstring	extract chinese from file or folder.
			e.g. trans getstring path				
	update		update translation to dictionary.
			e.g. trans update chinese.txt translate.txt				
	translate	translate file or folder.
			e.g. trans translate src_path des_path
	
Remark: Supports .lua, .prefab, .tab file`)
}

func main() {
	//	defer profile.Start(profile.CPUProfile).Stop()
	//	defer profile.Start(profile.MemProfile).Stop()
	switch len(os.Args) {
	case 3:
		if strings.EqualFold(os.Args[1], "getstring") {
			functool.GetString(os.Args[2])
		} else {
			useage()
		}
	case 4:
		if strings.EqualFold(os.Args[1], "update") {
			functool.Update(os.Args[2], os.Args[3])
		} else if strings.EqualFold(os.Args[1], "translate") {
			functool.Translate(os.Args[2], os.Args[3], 1)
		} else {
			useage()
		}
	case 5:
		if strings.EqualFold(os.Args[1], "translate") {
			queue, _ := strconv.ParseInt(os.Args[4], 10, 0)
			functool.Translate(os.Args[2], os.Args[3], int(queue))
		} else {
			useage()
		}
	default:
		useage()
	}
}
