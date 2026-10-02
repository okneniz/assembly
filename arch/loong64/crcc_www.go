package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWWW - crcc.w.w.w rd, rj, rk (DJK): carryless CRC32 of rj (word) into
// rk's accumulator.
type CrccWWW struct {
	rd, rj, rk uint8
}

func (i CrccWWW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.w.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWWW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crcc.w.w.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
