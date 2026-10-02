package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SraD - sra.d rd, rj, rk (3R): rd = rj >>a (rk & 63).
type SraD struct {
	base

	rd, rj, rk uint8
}

func (i SraD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sra.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SraD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sra.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
