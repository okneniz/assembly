package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbWu - ammax_db.wu rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned. A built-in barrier.
type AmmaxDbWu struct {
	rd, rk, rj uint8
}

func (i AmmaxDbWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax_db.wu"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
