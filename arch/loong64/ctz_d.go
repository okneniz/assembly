package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtzD - ctz.d rd, rj (2R): rd = the count of trailing zeros of rj.
type CtzD struct {
	base

	rd, rj uint8
}

func (i CtzD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ctz.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtzD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ctz.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
