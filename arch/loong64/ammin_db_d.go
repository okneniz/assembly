package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminDbD - ammin_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), signed. A built-in barrier.
type AmminDbD struct {
	rd, rk, rj uint8
}

func (i AmminDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammin_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmminDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
