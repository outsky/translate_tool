package dic

import (
	"bytes"
	"fmt"
	"trans/filetool"
	"trans/log"
)

type dic struct {
	name    string
	line    [][]byte
	trans   map[string]string
	key2idx map[string]int
}

func New(file string) *dic {
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
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
		return ins
	}
	for i := 1; i < len(all); i++ {
		v := all[i]
		linev := bytes.Split(v, []byte{0x09})
		if len(linev) != 2 || len(linev[0]) == 0 || len(linev[1]) == 0 {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic abnormal] file:%s, line:%d, data:%s", file, i+1, v))
			continue
		}
		key := string(linev[0])
		if _, ok := ins.trans[key]; ok {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic repeat] file:%s, line:%d, data:%s", file, i+1, key))
			continue
		}
		ins.trans[key] = string(linev[1])
		ins.line = append(ins.line, v)
		ins.key2idx[key] = len(ins.line) - 1
	}
	return ins
}

func NewOnly(file string) *dic {
	return &dic{
		name:    file,
		trans:   make(map[string]string),
		key2idx: make(map[string]int),
	}
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
		elemv := bytes.Split(elem, []byte{0x09})
		if len(elemv) != 2 || len(elemv[0]) == 0 || len(elemv[1]) == 0 {
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
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
	}
}
