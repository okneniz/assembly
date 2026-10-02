package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtW - stgt.w rd, rj, rk (DJK): store the low word of rd with a bounds
// check against rk, trapping outside.
type StgtW struct {
	rd, rj, rk uint8
}

func (i StgtW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stgt.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
