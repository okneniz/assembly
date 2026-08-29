package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lddir - lddir rd, rj, ui8: read a page-walk directory entry;
// rd = the descriptor at level ui8 for address rj.
type Lddir struct {
	base

	rd, rj uint8
	imm    imm
}

// NewLddir - lddir rd, rj, ui8.
func NewLddir(rd, rj Reg, v UImm8) Instr {
	return Lddir{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeLddir(w uint32, addr uint64) Instr {
	return Lddir{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 8))),
	}
}

func (i Lddir) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lddir %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Lddir) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["lddir"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 8)

	return writeWord(w, word)
}

func (i Lddir) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lddir",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
