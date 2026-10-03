package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cpucfg - cpucfg rd, rj (2R): rd = the configuration register selected by rj.
type Cpucfg struct {
	rd, rj uint8
}

func (i Cpucfg) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["cpucfg"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Cpucfg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cpucfg %s, %s", laRegName(i.rd), laRegName(i.rj))
}
