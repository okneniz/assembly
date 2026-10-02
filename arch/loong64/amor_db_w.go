package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmorDbW - amor_db.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] |= rk. A built-in barrier.
type AmorDbW struct {
	base

	rd, rk, rj uint8
}

func (i AmorDbW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amor_db.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmorDbW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amor_db.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
