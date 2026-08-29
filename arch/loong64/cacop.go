package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cacop - cacop op, rj, si12: the cache operation op on the block at
// rj + si12 (an unscaled byte offset).
type Cacop struct {
	base

	op  imm
	rj  uint8
	off imm
}

// NewCacop - cacop op, rj, si12 (the assembly operand order).
func NewCacop(op UImm5, rj Reg, off Imm12) Instr {
	return Cacop{
		op:  immNum(op.Val()),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

func decodeCacop(w uint32, addr uint64) Instr {
	return Cacop{
		base: newBase(addr, w),
		op:   immNum(int64(uField(w, 0, 5))),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func (i Cacop) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cacop %s, %s, %s", i.op.text(), laRegName(i.rj), i.off.text())
}

func (i Cacop) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["cacop"][0] |
		scatterU(i.op.val, 0, 5) | uint32(i.rj)<<5 | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}

func (i Cacop) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"cacop",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"op": i.op.val, "rj": laRegName(i.rj), "off": i.off.val},
	)
}
