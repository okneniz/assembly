package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxD - ldx.d rd, rj, rk (3R): rd = MEM[rj + rk].
type LdxD struct {
	base

	rd, rj, rk uint8
}

func (i LdxD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldx.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
