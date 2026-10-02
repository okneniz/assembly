package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScQ - sc.q rd, rk, rj (3R): store the {rd, rd+1} pair (16 bytes) to [rk]; rj - a hint (0 = none).
type ScQ struct {
	rd, rk, rj uint8
}

func (i ScQ) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sc.q %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i ScQ) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sc.q"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
