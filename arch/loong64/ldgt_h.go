package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdgtH - ldgt.h rd, rj, rk (DJK): load a half with a bounds check against
// rk, trapping outside.
type LdgtH struct {
	base

	rd, rj, rk uint8
}

func (i LdgtH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldgt.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdgtH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldgt.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
