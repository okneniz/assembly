package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Masknez - masknez rd, rj, rk (3R): rd = (rk != 0) ? rj : 0.
type Masknez struct {
	base

	rd, rj, rk uint8
}

func (i Masknez) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("masknez %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Masknez) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["masknez"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
