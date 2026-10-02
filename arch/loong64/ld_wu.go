package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdWu - ld.wu rd, rj, si12 (2RI12): rd = zero32(MEM[rj + si12]) (the
// offset is an unscaled byte offset).
type LdWu struct {
	rd, rj uint8
	imm    imm
}

func (i LdWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i LdWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ld.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
