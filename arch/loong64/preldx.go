package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Preldx - preldx hint, rj, rk (Ud5JK): preload MEM[rj + rk] into the
// cache (hint selects the operation; the manual prints the hint first).
type Preldx struct {
	rj, rk uint8
	hint   imm
}

func (i Preldx) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["preldx"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10 | scatterU(i.hint.val, 0, 5)

	return writeWord(w, word)
}

func (i Preldx) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("preldx %s, %s, %s", i.hint.text(), laRegName(i.rj), laRegName(i.rk))
}
