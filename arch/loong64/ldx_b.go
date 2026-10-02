package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxB - ldx.b rd, rj, rk (3R): rd = sign8(MEM[rj + rk]).
type LdxB struct {
	base

	rd, rj, rk uint8
}

func (i LdxB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldx.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
