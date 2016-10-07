// Process .ini files
// Only deal with values, leave comments, headers, keys unchanged
package ini

import ()

type ini struct {
	filename string
}

func New(file string) *ini {
	return &ini{file}
}

func shouldRecord(txt []byte) bool {
	for _, c := range txt {
		if c&0x80 != 0 {
			return true
		}
	}
	return false
}

func (t *ini) GetString(text []byte) ([][]byte, []int, []int, error) {
	var entryStart []int
	var entryEnd []int
	var entryTotal [][]byte
	frecord := func(nStart, nEnd int) {
		slice := text[nStart:nEnd]
		if shouldRecord(slice) {
			entryStart = append(entryStart, nStart)
			entryEnd = append(entryEnd, nEnd)
			entryTotal = append(entryTotal, slice)
		}
	}

	const (
		normal = iota
		dealValue
		dealComment
	)

	state := normal

	nStart := 0
	length := len(text)
	for i := 0; i < length; i++ {
		c := text[i]
		switch state {
		case normal:
			if c == ';' {
				state = dealComment
			} else if c == '=' {
				nStart = i + 1
				state = dealValue
			}
		case dealValue:
			if c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == ';' || i >= (length-1) {
				end := i
				if i >= (length - 1) {
					end = length
				}
				frecord(nStart, end)
				state = normal
			}
		case dealComment:
			if c == '\r' || c == '\n' {
				state = normal
			}
		}
	}
	return entryTotal, entryStart, entryEnd, nil
}

func (t *ini) Pretreat(trans []byte) []byte {
	return trans
}
