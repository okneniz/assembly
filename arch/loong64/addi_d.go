package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddiD - addi.d rd, rj, si12 (2RI12): rd = rj + si12.
type AddiD struct {
	rd, rj uint8
	imm    imm
}

func (i AddiD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["addi.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i AddiD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addi.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
