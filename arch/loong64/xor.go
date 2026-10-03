package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xor - xor rd, rj, rk (3R): rd = rj ^ rk.
type Xor struct {
	rd, rj, rk uint8
}

func (i Xor) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["xor"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Xor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("xor %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
