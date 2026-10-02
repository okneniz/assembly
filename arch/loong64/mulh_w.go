package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulhW - mulh.w rd, rj, rk (3R): rd = high half of the signed 64-bit product low32(rj) * low32(rk).
type MulhW struct {
	rd, rj, rk uint8
}

func (i MulhW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulhW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mulh.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
