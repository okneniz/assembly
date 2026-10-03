package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmxorD - amxor.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] ^= rk.
type AmxorD struct {
	rd, rk, rj uint8
}

func (i AmxorD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amxor.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmxorD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amxor.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
