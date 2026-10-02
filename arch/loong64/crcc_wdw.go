package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWDW - crcc.w.d.w rd, rj, rk (DJK): carryless CRC32 of rj (double) into
// rk's accumulator.
type CrccWDW struct {
	rd, rj, rk uint8
}

func (i CrccWDW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.d.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWDW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crcc.w.d.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
