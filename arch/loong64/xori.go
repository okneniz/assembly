package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xori - xori rd, rj, ui12 (2RI12): rd = rj ^ ui12.
type Xori struct {
	base

	rd, rj uint8
	imm    imm
}

func (i Xori) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("xori %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Xori) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["xori"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 12)

	return writeWord(w, word)
}
