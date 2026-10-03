package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// And - and rd, rj, rk (3R): rd = rj & rk.
type And struct {
	rd, rj, rk uint8
}

func (i And) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["and"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i And) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("and %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
