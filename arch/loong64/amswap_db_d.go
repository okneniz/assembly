package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmswapDbD - amswap_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = rk. A built-in barrier.
type AmswapDbD struct {
	base

	rd, rk, rj uint8
}

func (i AmswapDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amswap_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmswapDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amswap_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
