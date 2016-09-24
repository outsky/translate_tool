package common

import (
	"bytes"
)

var (
	cr byte = 0x0d //回车CR
	lf byte = 0x0a //换行LF
)

const (
	state_normal  = iota //正常状态
	state_working        //工作状态
)

type common struct{}

func New() *common {
	return &common{}
}

func (t *common) GetString(text []byte) ([][]byte, []int, []int, error) {
	var entryStart []int
	var entryEnd []int
	var entryTotal [][]byte
	frecord := func(nStart, nEnd int) {
		slice := text[nStart:nEnd]
		if len(bytes.TrimSpace(slice)) > 0 {
			entryStart = append(entryStart, nStart)
			entryEnd = append(entryEnd, nEnd)
			entryTotal = append(entryTotal, slice)
		}
	}
	isch := func(b byte) bool {
		return b&0x80 != 0
	}
	nState := state_normal
	nStart := 0
	length := len(text)
	for i := 0; i < length; i++ {
		switch nState {
		case state_normal:
			if isch(text[i]) {
				nStart = i
				nState = state_working
			}
		case state_working:
			if !isch(text[i]) {
				frecord(nStart, i)
				nState = state_normal
			}
		}
	}
	return entryTotal, entryStart, entryEnd, nil
}

func (t *common) Pretreat(trans []byte) []byte {
	return trans
}
