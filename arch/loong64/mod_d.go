package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModD - mod.d rd, rj, rk (3R): rd = rj % rk (signed).
type ModD struct {
	rd, rj, rk uint8
}

func (i ModD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i ModD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mod.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
