package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdgtD - ldgt.d rd, rj, rk (DJK): load a double word with a bounds check
// against rk, trapping outside.
type LdgtD struct {
	rd, rj, rk uint8
}

func (i LdgtD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldgt.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdgtD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldgt.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
