package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lu52iD - lu52i.d rd, rj, si12 (2RI12): rd = (rj & low52) | (si12 << 52).
type Lu52iD struct {
	rd, rj uint8
	imm    imm
}

func (i Lu52iD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["lu52i.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Lu52iD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lu52i.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
