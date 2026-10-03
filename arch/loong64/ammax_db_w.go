package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbW - ammax_db.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed. A built-in barrier.
type AmmaxDbW struct {
	rd, rk, rj uint8
}

func (i AmmaxDbW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax_db.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmmaxDbW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
