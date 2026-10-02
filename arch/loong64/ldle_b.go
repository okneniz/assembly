package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdleB - ldle.b rd, rj, rk (DJK): load a byte with a bounds check against
// rk, trapping outside.
type LdleB struct {
	base

	rd, rj, rk uint8
}

func (i LdleB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldle.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdleB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ldle.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
