package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulwDWu - mulw.d.wu rd, rj, rk (3R): rd = the unsigned 64-bit product low32(rj) * low32(rk).
type MulwDWu struct {
	rd, rj, rk uint8
}

func (i MulwDWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulw.d.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulwDWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mulw.d.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
