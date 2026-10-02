package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxW - ammax.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed.
type AmmaxW struct {
	base

	rd, rk, rj uint8
}

func (i AmmaxW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
