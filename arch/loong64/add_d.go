package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddD - add.d rd, rj, rk (3R): rd = rj + rk.
type AddD struct {
	rd, rj, rk uint8
}

func (i AddD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i AddD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["add.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
