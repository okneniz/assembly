package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RotriD - rotri.d rd, rj, ui6 (2RI6): rd = rotate_right64(rj, ui6 & 63).
type RotriD struct {
	rd, rj uint8
	imm    imm
}

func (i RotriD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["rotri.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 6)

	return writeWord(w, word)
}

func (i RotriD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rotri.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
