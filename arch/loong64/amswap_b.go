package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapB - amswap.b rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk.
type AmswapB struct {
	rd, rk, rj uint8
}

func (i AmswapB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amswap.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
