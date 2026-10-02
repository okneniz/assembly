package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdH - ld.h rd, rj, si12 (2RI12): rd = sign16(MEM[rj + si12]) (the
// offset is an unscaled byte offset).
type LdH struct {
	base

	rd, rj uint8
	imm    imm
}

func (i LdH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i LdH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ld.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
