package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminW - ammin.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), signed.
type AmminW struct {
	rd, rk, rj uint8
}

func (i AmminW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammin.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
