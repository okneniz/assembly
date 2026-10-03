package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StleH - stle.h rd, rj, rk (DJK): store the low half of rd with a bounds
// check against rk, trapping outside.
type StleH struct {
	rd, rj, rk uint8
}

func (i StleH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stle.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StleH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stle.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
