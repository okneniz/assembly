package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Orn - orn rd, rj, rk (3R): rd = rj | ~rk.
type Orn struct {
	base

	rd, rj, rk uint8
}

func (i Orn) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("orn %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Orn) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["orn"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
