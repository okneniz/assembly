package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CloW - clo.w rd, rj (2R): rd = the count of leading ones of low32(rj).
type CloW struct {
	rd, rj uint8
}

func (i CloW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clo.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CloW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["clo.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
