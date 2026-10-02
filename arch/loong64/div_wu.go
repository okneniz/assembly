package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DivWu - div.wu rd, rj, rk (3R): rd = rj / rk (unsigned).
type DivWu struct {
	rd, rj, rk uint8
}

func (i DivWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i DivWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["div.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
