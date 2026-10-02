package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SllD - sll.d rd, rj, rk (3R): rd = rj << (rk & 63).
type SllD struct {
	base

	rd, rj, rk uint8
}

func (i SllD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sll.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SllD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sll.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
