package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lu32iD - lu32i.d rd, si20 (1RI20): rd = si20 << 32 concatenated into
// bits 51:32.
type Lu32iD struct {
	rd  uint8
	imm imm
}

func (i Lu32iD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["lu32i.d"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Lu32iD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lu32i.d %s, %s", laRegName(i.rd), i.imm.text())
}
