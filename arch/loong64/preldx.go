package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Preldx - preldx hint, rj, rk (Ud5JK): preload MEM[rj + rk] into the
// cache (hint selects the operation; the manual prints the hint first).
type Preldx struct {
	base

	rj, rk uint8
	hint   imm
}

// NewPreldx - preldx hint, rj, rk.
func NewPreldx(hint UImm5, rj, rk Reg) Instr {
	return Preldx{
		rj:   rj.Num(),
		rk:   rk.Num(),
		hint: immNum(hint.Val()),
	}
}

func decodePreldx(w uint32, addr uint64) Instr {
	return Preldx{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
	}
}

func (i Preldx) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("preldx %s, %s, %s", i.hint.text(), laRegName(i.rj), laRegName(i.rk))
}

func (i Preldx) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["preldx"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10 | scatterU(i.hint.val, 0, 5)

	return writeWord(w, word)
}

func (i Preldx) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"preldx",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"hint": i.hint.val, "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
