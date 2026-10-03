package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StxB - stx.b rd, rj, rk (3R): MEM[rj + rk] = low8(rd).
type StxB struct {
	rd, rj, rk uint8
}

func (i StxB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stx.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StxB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stx.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
