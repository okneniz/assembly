package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RdtimeD - rdtime.d rd, rj (2R): rd = the 64-bit stable counter + rj.
type RdtimeD struct {
	rd, rj uint8
}

func (i RdtimeD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rdtime.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RdtimeD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["rdtime.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
