package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StxH - stx.h rd, rj, rk (3R): MEM[rj + rk] = low16(rd).
type StxH struct {
	base

	rd, rj, rk uint8
}

func (i StxH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stx.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StxH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stx.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
