package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubW - sub.w rd, rj, rk (3R): rd = sign32(rj - rk).
type SubW struct {
	rd, rj, rk uint8
}

func (i SubW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sub.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SubW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sub.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
