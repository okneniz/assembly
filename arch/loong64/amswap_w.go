package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapW - amswap.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk.
type AmswapW struct {
	rd, rk, rj uint8
}

func (i AmswapW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amswap.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmswapW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
