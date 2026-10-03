package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulwDW - mulw.d.w rd, rj, rk (3R): rd = sign64(low32(rj) * low32(rk)).
type MulwDW struct {
	rd, rj, rk uint8
}

func (i MulwDW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mulw.d.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulwDW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulw.d.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
