package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdptrW - ldptr.w rd, rj, offs (DJSk14): rd = sign32(MEM[rj + offs])
// (offs is a word-scaled byte offset, stored raw).
type LdptrW struct {
	rd, rj uint8
	off    imm
}

func (i LdptrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldptr.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdptrW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldptr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}
