package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ori - ori rd, rj, ui12 (2RI12): rd = rj | ui12.
type Ori struct {
	rd, rj uint8
	imm    imm
}

func (i Ori) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ori"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Ori) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ori %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
