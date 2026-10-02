package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdW - ld.w rd, rj, si12 (2RI12): rd = sign32(MEM[rj + si12]) (the
// offset is an unscaled byte offset).
type LdW struct {
	base

	rd, rj uint8
	imm    imm
}

func (i LdW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i LdW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ld.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
