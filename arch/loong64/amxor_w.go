package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmxorW - amxor.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] ^= rk.
type AmxorW struct {
	rd, rk, rj uint8
}

func (i AmxorW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amxor.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmxorW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amxor.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}
