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

// NewPreld - preld hint, rj, si12.
func NewPreld(hint UImm5, rj Reg, off Imm12) Instr {
	return Preld{
		rj:   rj.Num(),
		hint: immNum(hint.Val()),
		off:  immNum(off.Val()),
	}
}

func decodePreld(w uint32, addr uint64) Instr {
	return Preld{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
		off:  immNum(sField(w, 10, 12)),
	}
}

func (i Preld) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("preld %s, %s, %s", i.hint.text(), laRegName(i.rj), i.off.text())
}

func (i Preld) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["preld"][0] |
		uint32(i.rj)<<5 | scatterU(i.hint.val, 0, 5) | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}

func (i Preld) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"preld",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"hint": i.hint.val, "rj": laRegName(i.rj), "off": i.off.val},
	)
}
