package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulhD - mulh.d rd, rj, rk (3R): rd = high half of the signed 128-bit product rj * rk.
type MulhD struct {
	rd, rj, rk uint8
}

func (i MulhD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulhD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mulh.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
