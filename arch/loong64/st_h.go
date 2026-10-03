package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StH - st.h rd, rj, si12 (2RI12): MEM[rj + si12] = low16(rd) (the
// offset is an unscaled byte offset).
type StH struct {
	rd, rj uint8
	imm    imm
}

func (i StH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["st.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i StH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("st.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
