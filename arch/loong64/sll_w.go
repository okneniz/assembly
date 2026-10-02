package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SllW - sll.w rd, rj, rk (3R): rd = sign32(rj << (rk & 31)).
type SllW struct {
	base

	rd, rj, rk uint8
}

func (i SllW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sll.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SllW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sll.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
