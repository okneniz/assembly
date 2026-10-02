package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Maskeqz - maskeqz rd, rj, rk (3R): rd = (rk == 0) ? rj : 0.
type Maskeqz struct {
	base

	rd, rj, rk uint8
}

func (i Maskeqz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("maskeqz %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Maskeqz) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["maskeqz"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
