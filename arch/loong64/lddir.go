package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lddir - lddir rd, rj, ui8: read a page-walk directory entry;
// rd = the descriptor at level ui8 for address rj.
type Lddir struct {
	rd, rj uint8
	imm    imm
}

func (i Lddir) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["lddir"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 8)

	return writeWord(w, word)
}

func (i Lddir) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lddir %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}
