package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Or - or rd, rj, rk (3R): rd = rj | rk.
type Or struct {
	rd, rj, rk uint8
}

func (i Or) ObjDump(_ disasm.ViewCtx) string {
	if i.rk == 0 {
		return fmt.Sprintf("move %s, %s", laRegName(i.rd), laRegName(i.rj))
	}

	return fmt.Sprintf("or %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Or) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["or"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
