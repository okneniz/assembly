package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevhD - revh.d rd, rj (2R): rd = rj with all four halfwords reversed.
type RevhD struct {
	rd, rj uint8
}

func (i RevhD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revh.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevhD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["revh.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
