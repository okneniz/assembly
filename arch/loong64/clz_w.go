package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ClzW - clz.w rd, rj (2R): rd = the count of leading zeros of low32(rj).
type ClzW struct {
	rd, rj uint8
}

func (i ClzW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clz.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ClzW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["clz.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
