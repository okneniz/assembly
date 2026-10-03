package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmandDbD - amand_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] &= rk. A built-in barrier.
type AmandDbD struct {
	rd, rk, rj uint8
}

func (i AmandDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amand_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmandDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amand_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
