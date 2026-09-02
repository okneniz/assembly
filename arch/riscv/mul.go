package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mul - mul rd, rs1, rs2.
type Mul struct {
	base

	rd, rs1, rs2 string
}

// Mul - mul rd, rs1, rs2.
func (Builder) Mul(rd, rs1, rs2 Reg) Instr {
	return Mul{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeMul(w uint32, addr uint64) Instr {
	return Mul{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Mul) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mul %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Mul) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["mul"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mul) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mul",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newMul(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mul")
	if err != nil {
		return nil, err
	}

	return Mul{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
