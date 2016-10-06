package dic

import (
	"bytes"
	"fmt"
	"strings"

	"trans/filetool"
	"trans/log"
)

type dic struct {
	name    string
	line    [][]byte
	trans   map[string]string
	key2idx map[string]int
}

func NewDic(file string) *dic {
	ins := &dic{
		name:    file,
		trans:   make(map[string]string),
		key2idx: make(map[string]int),
	}
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(file, "utf8")
	defer ft.SetEncoding(file, oldEncode)
	all, err := ft.ReadFileLine(file)
	if err != nil {
		log.Info("fc", err)
		return ins
	}
	for i := 1; i < len(all); i++ {
		v := all[i]
		linev := bytes.Split(v, []byte{'\t'})
		// ID	File	Original	Translation
		if len(linev) < 4 || len(linev[2]) == 0 || len(linev[3]) == 0 {
			log.Error("fc", fmt.Sprintf("[dic abnormal] file:%s, line:%d", file, i+1))
			continue
		}
		key := string(linev[2])
		if _, ok := ins.trans[key]; ok {
			log.Error("fc", fmt.Sprintf("[dic repeat] file:%s, line:%d, key:%s", file, i+1, key))
			continue
		}
		ins.trans[key] = string(linev[3])
		ins.line = append(ins.line, v)
		ins.key2idx[key] = len(ins.line) - 1
	}
	return ins
}

func (d *dic) Query(text []byte) ([]byte, bool) {
	stext := string(text)
	strans, ok := d.trans[stext]
	return []byte(strans), ok
}

func (d *dic) Append(text []byte, trans []byte) {
	stext := string(text)
	strans := string(trans)
	if _, ok := d.trans[stext]; ok {
		d.trans[stext] = strans
		d.line[d.key2idx[stext]] = []byte(fmt.Sprintf("%s\t%s", stext, strans))
	} else {
		d.trans[stext] = strans
		d.line = append(d.line, []byte(fmt.Sprintf("%s\t%s", stext, strans)))
		d.key2idx[stext] = len(d.line) - 1
	}
}

func (d *dic) GetLine() ([][]byte, [][]byte) {
	var text [][]byte
	var trans [][]byte
	for i := 0; i < len(d.line); i++ {
		elem := d.line[i]
		elemv := bytes.Split(elem, []byte{'\t'})
		if len(elemv) < 2 || len(elemv[0]) == 0 || len(elemv[1]) == 0 {
			continue
		}
		text = append(text, elemv[0])
		trans = append(trans, elemv[1])
	}
	return text, trans
}

func (d *dic) Save() {
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(d.name, "utf8")
	defer ft.SetEncoding(d.name, oldEncode)
	var all [][]byte
	all = append(all, []byte("Original\tTranslation"))
	all = append(all, d.line...)
	err := ft.SaveFileLine(d.name, all)
	if err != nil {
		log.Error("fc", err)
	}
}

type upt struct {
	name    string
	line    [][]byte
	trans   map[string]string
	key2idx map[string]int
}

func NewUpt(file string) *upt {
	return &upt{
		name:    file,
		trans:   make(map[string]string),
		key2idx: make(map[string]int),
	}
}

func (d *upt) Append(text []byte, trans []byte, source string) bool {
	stext := string(text)
	strans := string(trans)
	if _, ok := d.trans[stext]; ok {
		d.trans[stext] = strans
		contextv := bytes.Split(d.line[d.key2idx[stext]], []byte("\t"))
		context := string(contextv[2])
		subcontextv := strings.Split(context, ";")
		bExist := false
		for _, v := range subcontextv {
			if strings.EqualFold(source, v) {
				bExist = true
				break
			}
		}
		if !bExist {
			d.line[d.key2idx[stext]] = []byte(fmt.Sprintf("%s\t%s\t%s", stext, strans, strings.Join([]string{source, context}, ";")))
		} else {
			d.line[d.key2idx[stext]] = []byte(fmt.Sprintf("%s\t%s\t%s", stext, strans, source))
		}
		return false
	} else {
		d.trans[stext] = strans
		d.line = append(d.line, []byte(fmt.Sprintf("%s\t%s\t%s", stext, strans, source)))
		d.key2idx[stext] = len(d.line) - 1
		return true
	}
}

func (d *upt) Save() {
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(d.name, "utf8")
	defer ft.SetEncoding(d.name, oldEncode)
	var all [][]byte
	all = append(all, []byte("Original\tTranslation\tSource"))
	all = append(all, d.line...)
	err := ft.SaveFileLine(d.name, all)
	if err != nil {
		log.Error("fc", err)
	}
}
