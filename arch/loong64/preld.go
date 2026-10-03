package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Preld - preld hint, rj, si12 (Ud5JSk12): preload MEM[rj + si12] into
// the cache (hint selects the operation; the manual prints the hint
// first). The offset is an unscaled byte offset.
type Preld struct {
	rj   uint8
	hint imm
	off  imm
}

func (i Preld) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["preld"][0] |
		uint32(i.rj)<<5 | scatterU(i.hint.val, 0, 5) | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}

func (i Preld) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("preld %s, %s, %s", i.hint.text(), laRegName(i.rj), i.off.text())
}
