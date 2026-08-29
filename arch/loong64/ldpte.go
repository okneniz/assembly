package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldpte - ldpte rj, ui8: load the page table entry for address rj,
// level ui8.
type Ldpte struct {
	base

	rj  uint8
	imm imm
}

// NewLdpte - ldpte rj, ui8.
func NewLdpte(rj Reg, v UImm8) Instr {
	return Ldpte{
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeLdpte(w uint32, addr uint64) Instr {
	return Ldpte{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 8))),
	}
}

func (i Ldpte) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldpte %s, %s", laRegName(i.rj), i.imm.text())
}

func (i Ldpte) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldpte"][0] |
		uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 8)

	return writeWord(w, word)
}

func (i Ldpte) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldpte",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
