package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lu12iW - lu12i.w rd, si20 (1RI20): rd = si20 << 12 (the low 32 bits are
// emptied).
type Lu12iW struct {
	rd  uint8
	imm imm
}

func (i Lu12iW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["lu12i.w"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Lu12iW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lu12i.w %s, %s", laRegName(i.rd), i.imm.text())
}
