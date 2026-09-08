package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lb - lb rd, off(rs1).
type Lb struct {
	base

	rd, rs1 string
	off     imm
}

// Lb - lb rd, off(rs1).
func (Builder) Lb(rd, rs1 Reg, off Off) Instr {
	return Lb{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeLb(w uint32) Instr {
	return Lb{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Lb) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lb %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Lb) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lb"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func newLb(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lb", false)
	if err != nil {
		return nil, err
	}

	return Lb{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
