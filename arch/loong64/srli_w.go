package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrliW - srli.w rd, rj, ui5 (2RI5): rd = sign32(rj >> ui5) logical.
type SrliW struct {
	rd, rj uint8
	imm    imm
}

func (i SrliW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["srli.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 5)

	return writeWord(w, word)
}

func (i SrliW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srli.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
