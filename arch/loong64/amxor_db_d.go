package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmxorDbD - amxor_db.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] ^= rk. A built-in barrier.
type AmxorDbD struct {
	rd, rk, rj uint8
}

func (i AmxorDbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amxor_db.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmxorDbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amxor_db.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
