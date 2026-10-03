package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StB - st.b rd, rj, si12 (2RI12): MEM[rj + si12] = low8(rd) (the
// offset is an unscaled byte offset).
type StB struct {
	rd, rj uint8
	imm    imm
}

func (i StB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["st.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i StB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("st.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
