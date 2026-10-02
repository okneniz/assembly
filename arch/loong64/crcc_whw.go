package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWHW - crcc.w.h.w rd, rj, rk (DJK): carryless CRC32 of rj (half) into
// rk's accumulator.
type CrccWHW struct {
	base

	rd, rj, rk uint8
}

func (i CrccWHW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.h.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWHW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crcc.w.h.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
