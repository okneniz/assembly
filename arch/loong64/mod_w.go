package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModW - mod.w rd, rj, rk (3R): rd = rj % rk (signed).
type ModW struct {
	rd, rj, rk uint8
}

func (i ModW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mod.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i ModW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
