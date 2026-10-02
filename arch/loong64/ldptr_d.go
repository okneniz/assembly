package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdptrD - ldptr.d rd, rj, offs (DJSk14): rd = MEM[rj + offs] (offs is
// a word-scaled byte offset, stored raw).
type LdptrD struct {
	rd, rj uint8
	off    imm
}

func (i LdptrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldptr.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdptrD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldptr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}
