package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sltu - sltu rd, rj, rk (3R): rd = (rj < rk) ? 1 : 0 (unsigned).
type Sltu struct {
	rd, rj, rk uint8
}

func (i Sltu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sltu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Sltu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sltu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
