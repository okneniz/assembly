package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LlD - ll.d rd, rj, offs (DJSk14): rd = MEM[rj + offs] with a reservation
// (load-linked 64 bits; offs - a byte offset, a multiple of 4, in +-16380).
type LlD struct {
	rd, rj uint8
	off    imm
}

func (i LlD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ll.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LlD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ll.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}
