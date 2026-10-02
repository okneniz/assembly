package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CloD - clo.d rd, rj (2R): rd = the count of leading ones of rj.
type CloD struct {
	rd, rj uint8
}

func (i CloD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clo.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CloD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["clo.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
