package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Invtlb - invtlb op, rj, rk: invalidate the TLB entries selected by
// op using rj and rk.
type Invtlb struct {
	base

	rj, rk uint8
	op     imm
}

// Invtlb - invtlb op, rj, rk (the assembly operand order).
func (Builder) Invtlb(op UImm5, rj, rk Reg) Instr {
	return Invtlb{
		rj: rj.Num(),
		rk: rk.Num(),
		op: immNum(op.Val()),
	}
}

func decodeInvtlb(w uint32, addr uint64) Instr {
	return Invtlb{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		op:   immNum(int64(uField(w, 0, 5))),
	}
}

func (i Invtlb) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("invtlb %s, %s, %s", i.op.text(), laRegName(i.rj), laRegName(i.rk))
}

func (i Invtlb) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["invtlb"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10 | scatterU(i.op.val, 0, 5)

	return writeWord(w, word)
}

func (i Invtlb) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"invtlb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"op": i.op.val, "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
