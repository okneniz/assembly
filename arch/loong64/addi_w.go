package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddiW - addi.w rd, rj, si12 (2RI12): rd = sign32(rj + si12).
type AddiW struct {
	rd, rj uint8
	imm    imm
}

func (i AddiW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addi.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i AddiW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["addi.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
