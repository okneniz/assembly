package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdgtW - ldgt.w rd, rj, rk (DJK): load a word with a bounds check against
// rk, trapping outside.
type LdgtW struct {
	rd, rj, rk uint8
}

func (i LdgtW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldgt.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdgtW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldgt.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
