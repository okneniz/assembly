package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasW - amcas.w rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj].
type AmcasW struct {
	rd, rk, rj uint8
}

func (i AmcasW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amcas.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmcasW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
