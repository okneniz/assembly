package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevbitD - bitrev.d rd, rj (2R): rd = rj with all 64 bits reversed.
type RevbitD struct {
	rd, rj uint8
}

func (i RevbitD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bitrev.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RevbitD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}
