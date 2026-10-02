package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StxW - stx.w rd, rj, rk (3R): MEM[rj + rk] = low32(rd).
type StxW struct {
	base

	rd, rj, rk uint8
}

func (i StxW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stx.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StxW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stx.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
