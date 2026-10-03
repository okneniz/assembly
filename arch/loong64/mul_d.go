package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulD - mul.d rd, rj, rk (3R): rd = the low 64 bits of the product rj * rk.
type MulD struct {
	rd, rj, rk uint8
}

func (i MulD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mul.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mul.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
