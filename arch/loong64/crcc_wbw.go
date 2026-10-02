package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWBW - crcc.w.b.w rd, rj, rk (DJK): carryless CRC32 of rj (byte) into
// rk's accumulator.
type CrccWBW struct {
	rd, rj, rk uint8
}

func (i CrccWBW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.b.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWBW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crcc.w.b.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
