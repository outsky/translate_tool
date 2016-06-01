package analysis

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type analysis struct {
	uc2hz map[string]string
}

var (
	ap byte = 0x27 //单引号'
	dq byte = 0x22 //双引号"
	sl byte = 0x5c //转义斜杠\\
	bs byte = 0x2d //横杠-
	bl byte = 0x5b //左中括号[
	br byte = 0x5d //右中括号]
	cr byte = 0x0d //回车CR
	lf byte = 0x0a //换行LF
	eq byte = 0x3d //等于号=
	sp byte = 0x20 //空格
	tb byte = 0x09 //tab制表符
	uu byte = 0x75 //u字符
)

const (
	state_normal          = iota //正常状态
	state_note_line              //注释一行
	state_note_section           //注释段落
	state_apostrophe             //'单引号'字符串
	state_double_quotes          //"双引号"字符串
	state_double_brackets        //[[中括号]]字符串
)

func New() *analysis {
	return &analysis{
		uc2hz: make(map[string]string),
	}
}

func (a *analysis) Uc2hanzi(uc string) (string, error) {
	if hz, ok := a.uc2hz[uc]; ok {
		return hz, nil
	}
	val2int, err := strconv.ParseInt(uc, 16, 32)
	if err != nil {
		return uc, err
	}
	context := fmt.Sprintf("%c", val2int)
	a.uc2hz[uc] = context
	return context, nil
}

func (a *analysis) GetRule(file string) (
	func(*[]byte) (*[][]byte, error),
	func(*[]byte, []byte, []byte) error,
	error) {
	filev := strings.Split(file, ".")
	file_ex := filev[len(filev)-1]
	switch file_ex {
	case "lua":
		return a.analysis_lua, a.translate_lua, nil
	case "prefab":
		return a.analysis_prefab, a.translate_prefab, nil
	case "tab":
		return a.analysis_tab, a.translate_tab, nil
	default:
		return nil, nil, errors.New(fmt.Sprintf("[file not rule] %s", file))
	}
}

func (a *analysis) analysis_lua(text *[]byte) (*[][]byte, error) {
	var cnEntry [][]byte
	frecord := func(start, end int) {
		bIsChinese := false
		for i := start; i <= end; i++ {
			if (*text)[i]&0x80 != 0 {
				bIsChinese = true
				break
			}
		}
		if bIsChinese {
			cnEntry = append(cnEntry, (*text)[start:end+1])
		}
	}
	nState := state_normal
	nStateStart := 0
	nSize := len(*text)
	for i := 0; i < nSize; i++ {
		if i+1 < nSize && (*text)[i] == sl &&
			((*text)[i+1] == ap || (*text)[i+1] == dq || (*text)[i+1] == sl) {
			i++
			continue
		}
		switch nState {
		case state_normal:
			switch (*text)[i] {
			case bs:
				if i+1 < nSize && (*text)[i+1] == bs {
					if i+3 < nSize && (*text)[i+2] == bl {
						nPos := i + 3
						for nPos < nSize && (*text)[nPos] == eq {
							if (*text)[nPos] == cr || (*text)[nPos] == lf {
								break
							}
							nPos++
						}
						if (*text)[nPos] == bl {
							i += (nPos - i)
							nState = state_note_section
						} else {
							i += 1
							nState = state_note_line
						}
					} else {
						i += 1
						nState = state_note_line
					}
				}
			case ap:
				nStateStart = i + 1
				nState = state_apostrophe
			case dq:
				nStateStart = i + 1
				nState = state_double_quotes
			case bl:
				if i+1 < nSize && (*text)[i+1] == bl {
					nStateStart = i + 2
					nState = state_double_brackets
					i += 1
				}
			}
		case state_note_line:
			if i+1 < nSize && (*text)[i] == cr && (*text)[i] == lf {
				i += 1
				nState = state_normal
			} else if (*text)[i] == cr || (*text)[i] == lf {
				nState = state_normal
			}
		case state_note_section:
			if (*text)[i] == br {
				nPos := i + 1
				for nPos < nSize && (*text)[nPos] == eq {
					if (*text)[nPos] == cr || (*text)[i] == lf {
						break
					}
					nPos++
				}
				if (*text)[nPos] == br {
					i += (nPos - i)
					nState = state_normal
				}
			}
		case state_apostrophe:
			if i+1 < nSize && (*text)[i] == cr && (*text)[i] == lf {
				frecord(nStateStart, i-1)
				i += 1
				nStateStart = i + 1
			} else if (*text)[i] == cr || (*text)[i] == lf {
				frecord(nStateStart, i-1)
				nStateStart = i + 1
			} else if (*text)[i] == ap {
				frecord(nStateStart, i-1)
				nState = state_normal
			}
		case state_double_quotes:
			if i+1 < nSize && (*text)[i] == cr && (*text)[i] == lf {
				frecord(nStateStart, i-1)
				i += 1
				nStateStart = i + 1
			} else if (*text)[i] == cr || (*text)[i] == lf {
				frecord(nStateStart, i-1)
				nStateStart = i + 1
			} else if (*text)[i] == dq {
				frecord(nStateStart, i-1)
				nState = state_normal
			}
		case state_double_brackets:
			if i+1 < nSize && (*text)[i] == cr && (*text)[i] == lf {
				frecord(nStateStart, i-1)
				i += 1
				nStateStart = i + 1
			} else if (*text)[i] == cr || (*text)[i] == lf {
				frecord(nStateStart, i-1)
				nStateStart = i + 1
			} else if (*text)[i] == br {
				if i+1 < nSize && (*text)[i+1] == br {
					frecord(nStateStart, i-1)
					i += 1
					nState = state_normal
				}
			}
		}
	}
	if nState != state_normal && nState != state_note_line {
		return &cnEntry, errors.New(fmt.Sprintf("%s state:%d", "file syntax error", nState))
	}
	return &cnEntry, nil
}

func (a *analysis) translate_lua(context *[]byte, sText []byte, trans []byte) error {
	(*context) = bytes.Replace(*context, sText, trans, 1)
	return nil
}

func (a *analysis) analysis_prefab(text *[]byte) (*[][]byte, error) {
	var cnEntry [][]byte
	tag := fmt.Sprintf("%c%c", sl, uu)
	frecord := func(start, end int) {
		unicode := string((*text)[start : end+1])
		index := strings.Index(unicode, tag)
		for ; index != -1; index = strings.Index(unicode, tag) {
			hanzi, err := a.Uc2hanzi(unicode[index+2 : index+6])
			if err != nil {
				panic(err)
			}
			unicode = strings.Replace(unicode, unicode[index:index+6], hanzi, 1)
		}
		cnEntry = append(cnEntry, []byte(unicode))
	}
	nState := state_normal
	nStateStart := 0
	nSize := len(*text)
	for i := 0; i < nSize; i++ {
		switch nState {
		case state_normal:
			switch (*text)[i] {
			case dq:
				nStateStart = i + 1
				nState = state_double_quotes
			}
		case state_double_quotes:
			if (*text)[i] == dq {
				frecord(nStateStart, i-1)
				nState = state_normal
			}
		}
	}
	if nState != state_normal {
		return &cnEntry, errors.New(fmt.Sprintf("%s state:%d", "file syntax error", nState))
	}
	return &cnEntry, nil
}

func (a *analysis) translate_prefab(context *[]byte, sText []byte, trans []byte) error {
	prefabformat := func(s string) string {
		length := len(s)
		for i := 0; i+5 < length; i++ {
			if s[i] == sl && s[i+1] == uu {
				upper := strings.ToUpper(s[i+2 : i+6])
				s = strings.Replace(s, s[i+2:i+6], upper, 1)
			}
		}
		return s
	}
	textQuoted := strconv.QuoteToASCII(string(sText))
	textUnquoted := prefabformat(textQuoted[1 : len(textQuoted)-1])
	textUnquoted = strings.Replace(textUnquoted, "\\\\", "\\", -1)
	transQuoted := strconv.QuoteToASCII(string(trans))
	transUnquoted := prefabformat(transQuoted[1 : len(transQuoted)-1])
	transUnquoted = strings.Replace(transUnquoted, "\\\\", "\\", -1)
	(*context) = bytes.Replace(*context, []byte(textUnquoted), []byte(transUnquoted), 1)
	return nil
}

func (a *analysis) analysis_tab(text *[]byte) (*[][]byte, error) {
	var cnEntry [][]byte
	frecord := func(nStart, nEnd int) {
		textv := bytes.Split((*text)[nStart:nEnd], []byte{tb})
		for _, v := range textv {
			v = bytes.Trim(v, string(sp))
			bIsChinese := false
			for i := 0; i < len(v); i++ {
				if v[i]&0x80 != 0 {
					bIsChinese = true
					break
				}
			}
			if bIsChinese {
				cnEntry = append(cnEntry, v)
			}
		}
	}
	nStart := 0
	length := len(*text)
	for i := 0; i < length; i++ {
		if (*text)[i] == cr || (*text)[i] == lf {
			frecord(nStart, i)
			nStart = i + 1
		} else if i+1 < length && (*text)[i] == cr && (*text)[i] == lf {
			frecord(nStart, i)
			nStart = i + 2
		}
	}
	return &cnEntry, nil
}

func (a *analysis) translate_tab(context *[]byte, sText []byte, trans []byte) error {
	(*context) = bytes.Replace(*context, sText, trans, 1)
	return nil
}
