package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revbit8B - bitrev.8b rd, rj (2R): rd = rj with the bits reversed in each byte.
type Revbit8B struct {
	base

	rd, rj uint8
}

func (i Revbit8B) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.8b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revbit8B) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bitrev.8b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
