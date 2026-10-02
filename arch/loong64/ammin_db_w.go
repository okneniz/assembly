package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminDbW - ammin_db.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), signed. A built-in barrier.
type AmminDbW struct {
	base

	rd, rk, rj uint8
}

func (i AmminDbW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin_db.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminDbW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammin_db.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
