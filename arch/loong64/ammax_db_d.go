package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbD - ammax_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed. A built-in barrier.
type AmmaxDbD struct {
	rd, rk, rj uint8
}

func (i AmmaxDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
