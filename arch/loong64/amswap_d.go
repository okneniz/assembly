package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapD - amswap.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk.
type AmswapD struct {
	rd, rk, rj uint8
}

func (i AmswapD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amswap.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
