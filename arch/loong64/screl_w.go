package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScrelW - screl.w rd, rj (DJ): store-release MEM[rj] = rd (32 bits).
type ScrelW struct {
	base

	rd, rj uint8
}

func (i ScrelW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("screl.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ScrelW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["screl.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
