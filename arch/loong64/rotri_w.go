package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RotriW - rotri.w rd, rj, ui5 (2RI5): rd = sign32(rotate_right32(rj, ui5)).
type RotriW struct {
	rd, rj uint8
	imm    imm
}

func (i RotriW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["rotri.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 5)

	return writeWord(w, word)
}

func (i RotriW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rotri.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
