package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Andi - andi rd, rj, ui12 (2RI12): rd = rj & ui12.
type Andi struct {
	rd, rj uint8
	imm    imm
}

func (i Andi) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["andi"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Andi) ObjDump(_ disasm.ViewCtx) string {
	if i.rd == 0 && i.rj == 0 && i.imm.val == 0 {
		return "nop"
	}

	return fmt.Sprintf("andi %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
