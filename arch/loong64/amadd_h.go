package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmaddH - amadd.h rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] += rk.
type AmaddH struct {
	rd, rk, rj uint8
}

func (i AmaddH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amadd.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmaddH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amadd.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
