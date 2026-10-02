package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasD - amcas.d rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj].
type AmcasD struct {
	rd, rk, rj uint8
}

func (i AmcasD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amcas.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
