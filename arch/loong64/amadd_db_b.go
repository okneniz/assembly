package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmaddDbB - amadd_db.b rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] += rk. A built-in barrier.
type AmaddDbB struct {
	rd, rk, rj uint8
}

func (i AmaddDbB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amadd_db.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmaddDbB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amadd_db.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
