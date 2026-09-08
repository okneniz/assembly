package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Flw - flw rd, off(rs1) (rd - floating).
type Flw struct {
	base

	rd, rs1 string
	off     imm
}

// Flw - flw rd, off(rs1); rd is an FP register taken by number
// (Reg 0..31 is printed ft0/fa0/...).
func (Builder) Flw(rd, rs1 Reg, off Off) Instr {
	return Flw{
		rd:  fpName(rd),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeFlw(w uint32) Instr {
	return Flw{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Flw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("flw %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Flw) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["flw"][0] | fregBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func newFlw(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "flw", true)
	if err != nil {
		return nil, err
	}

	return Flw{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
