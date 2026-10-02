package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revb4H - revb.4h rd, rj (2R): rd = rj with the bytes reversed in each halfword.
type Revb4H struct {
	base

	rd, rj uint8
}

func (i Revb4H) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.4h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revb4H) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["revb.4h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
