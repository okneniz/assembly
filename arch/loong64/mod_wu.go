package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModWu - mod.wu rd, rj, rk (3R): rd = rj % rk (unsigned).
type ModWu struct {
	rd, rj, rk uint8
}

func (i ModWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mod.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i ModWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
