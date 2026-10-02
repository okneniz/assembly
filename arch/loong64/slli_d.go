package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SlliD - slli.d rd, rj, ui6 (2RI6): rd = rj << (ui6 & 63).
type SlliD struct {
	rd, rj uint8
	imm    imm
}

func (i SlliD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slli.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i SlliD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["slli.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 6)

	return writeWord(w, word)
}
