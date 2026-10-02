package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevbD - revb.d rd, rj (2R): rd = rj with all eight bytes reversed.
type RevbD struct {
	base

	rd, rj uint8
}

func (i RevbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevbD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["revb.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
