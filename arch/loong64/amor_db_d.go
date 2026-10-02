package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmorDbD - amor_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] |= rk. A built-in barrier.
type AmorDbD struct {
	base

	rd, rk, rj uint8
}

func (i AmorDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amor_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmorDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amor_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
