package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtoW - cto.w rd, rj (2R): rd = the count of trailing ones of low32(rj).
type CtoW struct {
	rd, rj uint8
}

func (i CtoW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cto.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtoW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["cto.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
