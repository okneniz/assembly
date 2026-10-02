package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Nor - nor rd, rj, rk (3R): rd = ~(rj | rk).
type Nor struct {
	rd, rj, rk uint8
}

func (i Nor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("nor %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Nor) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["nor"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
