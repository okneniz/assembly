package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDu - ammax.du rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned.
type AmmaxDu struct {
	rd, rk, rj uint8
}

func (i AmmaxDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.du %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax.du"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
