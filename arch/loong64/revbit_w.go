package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevbitW - bitrev.w rd, rj (2R): rd = low32(rj) with all 32 bits reversed.
type RevbitW struct {
	rd, rj uint8
}

func (i RevbitW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevbitW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bitrev.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
