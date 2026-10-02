package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdD - ld.d rd, rj, si12 (2RI12): rd = MEM[rj + si12] (the offset is
// an unscaled byte offset).
type LdD struct {
	rd, rj uint8
	imm    imm
}

func (i LdD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i LdD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ld.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
