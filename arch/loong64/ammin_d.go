package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminD - ammin.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), signed.
type AmminD struct {
	rd, rk, rj uint8
}

func (i AmminD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammin.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
