package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Andn - andn rd, rj, rk (3R): rd = rj & ~rk.
type Andn struct {
	rd, rj, rk uint8
}

func (i Andn) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("andn %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Andn) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["andn"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
