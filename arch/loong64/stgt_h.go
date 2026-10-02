package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtH - stgt.h rd, rj, rk (DJK): store the low half of rd with a bounds
// check against rk, trapping outside.
type StgtH struct {
	base

	rd, rj, rk uint8
}

func (i StgtH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stgt.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
