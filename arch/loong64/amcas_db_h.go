package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasDbH - amcas_db.h rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj]. A built-in barrier.
type AmcasDbH struct {
	rd, rk, rj uint8
}

func (i AmcasDbH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amcas_db.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmcasDbH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas_db.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
