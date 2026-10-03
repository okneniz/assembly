package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxW - ldx.w rd, rj, rk (3R): rd = sign32(MEM[rj + rk]).
type LdxW struct {
	rd, rj, rk uint8
}

func (i LdxW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldx.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
