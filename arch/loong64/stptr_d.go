package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StptrD - stptr.d rd, rj, offs (DJSk14): MEM[rj + offs] = rd (offs is
// a word-scaled byte offset, stored raw).
type StptrD struct {
	rd, rj uint8
	off    imm
}

func (i StptrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stptr.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i StptrD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stptr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}
