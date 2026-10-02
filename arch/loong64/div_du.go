package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DivDu - div.du rd, rj, rk (3R): rd = rj / rk (unsigned).
type DivDu struct {
	rd, rj, rk uint8
}

func (i DivDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div.du %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i DivDu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["div.du"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
