package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StptrW - stptr.w rd, rj, offs (DJSk14): MEM[rj + offs] = low32(rd)
// (offs is a word-scaled byte offset, stored raw).
type StptrW struct {
	base

	rd, rj uint8
	off    imm
}

func (i StptrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stptr.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i StptrW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stptr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}
