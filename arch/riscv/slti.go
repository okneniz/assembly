package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slti - slti rd, rs1, imm.
type Slti struct {
	base

	rd, rs1 string
	imm     imm
}

// Slti - slti rd, rs1, imm.
func (Builder) Slti(rd, rs1 Reg, imm Imm12) Instr {
	return Slti{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

func decodeSlti(w uint32, addr uint64) Instr {
	return Slti{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func (i Slti) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slti %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Slti) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["slti"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|bits)
}

func (i Slti) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"slti",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newSlti(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "slti")
	if err != nil {
		return nil, err
	}

	return Slti{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
