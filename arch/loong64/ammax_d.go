package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxD - ammax.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed.
type AmmaxD struct {
	base

	rd, rk, rj uint8
}

func (i AmmaxD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
