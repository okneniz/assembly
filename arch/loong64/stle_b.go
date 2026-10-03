package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StleB - stle.b rd, rj, rk (DJK): store the low byte of rd with a bounds
// check against rk, trapping outside.
type StleB struct {
	rd, rj, rk uint8
}

func (i StleB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stle.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StleB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stle.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
