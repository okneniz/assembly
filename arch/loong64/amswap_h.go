package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapH - amswap.h rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk.
type AmswapH struct {
	rd, rk, rj uint8
}

func (i AmswapH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amswap.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmswapH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
