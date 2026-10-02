package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtoD - cto.d rd, rj (2R): rd = the count of trailing ones of rj.
type CtoD struct {
	base

	rd, rj uint8
}

func (i CtoD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cto.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtoD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["cto.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
