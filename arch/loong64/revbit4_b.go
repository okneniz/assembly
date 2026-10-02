package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revbit4B - bitrev.4b rd, rj (2R): rd = rj with the bits reversed in each nibble.
type Revbit4B struct {
	base

	rd, rj uint8
}

func (i Revbit4B) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.4b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revbit4B) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bitrev.4b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
