package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrlD - srl.d rd, rj, rk (3R): rd = rj >>l (rk & 63).
type SrlD struct {
	rd, rj, rk uint8
}

func (i SrlD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["srl.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SrlD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srl.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
