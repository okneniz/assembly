package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxWu - ldx.wu rd, rj, rk (3R): rd = zero32(MEM[rj + rk]).
type LdxWu struct {
	base

	rd, rj, rk uint8
}

func (i LdxWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldx.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
