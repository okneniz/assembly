package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ClzD - clz.d rd, rj (2R): rd = the count of leading zeros of rj.
type ClzD struct {
	rd, rj uint8
}

func (i ClzD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clz.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ClzD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["clz.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
