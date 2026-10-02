package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasB - amcas.b rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj].
type AmcasB struct {
	rd, rk, rj uint8
}

func (i AmcasB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amcas.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
