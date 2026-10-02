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
	base

	rj   uint8
	hint imm
	off  imm
}

func decodePreld(w uint32) Instr {
	return Preld{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
		off:  immNum(sField(w, 10, 12)),
	}
}

func (i Preld) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("preld %s, %s, %s", i.hint.text(), laRegName(i.rj), i.off.text())
}

func (i Preld) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["preld"][0] |
		uint32(i.rj)<<5 | scatterU(i.hint.val, 0, 5) | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}
