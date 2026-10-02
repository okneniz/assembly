package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StW - st.w rd, rj, si12 (2RI12): MEM[rj + si12] = low32(rd) (the
// offset is an unscaled byte offset).
type StW struct {
	base

	rd, rj uint8
	imm    imm
}

func (i StW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("st.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i StW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["st.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
