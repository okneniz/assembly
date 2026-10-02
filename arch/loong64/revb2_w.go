package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revb2W - revb.2w rd, rj (2R): rd = rj with the bytes reversed in each word.
type Revb2W struct {
	rd, rj uint8
}

func (i Revb2W) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.2w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revb2W) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["revb.2w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
