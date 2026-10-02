package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtD - stgt.d rd, rj, rk (DJK): store the double word of rd with a
// bounds check against rk, trapping outside.
type StgtD struct {
	base

	rd, rj, rk uint8
}

func (i StgtD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stgt.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
