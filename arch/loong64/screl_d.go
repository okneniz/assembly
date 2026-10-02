package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScrelD - screl.d rd, rj (DJ): store-release MEM[rj] = rd (64 bits).
type ScrelD struct {
	rd, rj uint8
}

func (i ScrelD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("screl.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ScrelD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["screl.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
