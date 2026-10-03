package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrlW - srl.w rd, rj, rk (3R): rd = sign32(rj >>l (rk & 31)).
type SrlW struct {
	rd, rj, rk uint8
}

func (i SrlW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["srl.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SrlW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srl.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
