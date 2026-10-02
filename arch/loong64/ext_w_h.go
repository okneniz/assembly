package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ExtWH - ext.w.h rd, rj (2R): rd = sign-extend of the low half of rj to the native width.
type ExtWH struct {
	rd, rj uint8
}

func (i ExtWH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ext.w.h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ExtWH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ext.w.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
