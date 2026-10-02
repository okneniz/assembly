package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbDu - ammax_db.du rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned. A built-in barrier.
type AmmaxDbDu struct {
	base

	rd, rk, rj uint8
}

func (i AmmaxDbDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.du %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbDu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax_db.du"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
