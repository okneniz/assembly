package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasDbB - amcas_db.b rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj]. A built-in barrier.
type AmcasDbB struct {
	rd, rk, rj uint8
}

func (i AmcasDbB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas_db.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasDbB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amcas_db.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
