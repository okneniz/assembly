package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DivD - div.d rd, rj, rk (3R): rd = rj / rk (signed).
type DivD struct {
	base

	rd, rj, rk uint8
}

func (i DivD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i DivD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["div.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
