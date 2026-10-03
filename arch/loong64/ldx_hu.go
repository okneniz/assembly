package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxHu - ldx.hu rd, rj, rk (3R): rd = the halfword at rj + rk, zero-extended.
type LdxHu struct {
	rd, rj, rk uint8
}

func (i LdxHu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldx.hu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxHu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.hu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
