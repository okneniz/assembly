package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrliD - srli.d rd, rj, ui6 (2RI6): rd = rj >> (ui6 & 63) logical.
type SrliD struct {
	rd, rj uint8
	imm    imm
}

func (i SrliD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["srli.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 6)

	return writeWord(w, word)
}

func (i SrliD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srli.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
