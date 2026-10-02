package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminWu - ammin.wu rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), unsigned.
type AmminWu struct {
	rd, rk, rj uint8
}

func (i AmminWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammin.wu"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
