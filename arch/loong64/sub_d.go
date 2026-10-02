package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubD - sub.d rd, rj, rk (3R): rd = rj - rk.
type SubD struct {
	rd, rj, rk uint8
}

func (i SubD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sub.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SubD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sub.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
