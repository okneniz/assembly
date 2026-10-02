package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldpte - ldpte rj, ui8: load the page table entry for address rj,
// level ui8.
type Ldpte struct {
	base

	rj  uint8
	imm imm
}

func (i Ldpte) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldpte %s, %s", laRegName(i.rj), i.imm.text())
}

func (i Ldpte) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldpte"][0] |
		uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 8)

	return writeWord(w, word)
}
