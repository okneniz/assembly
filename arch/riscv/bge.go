package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bge - bge rs1, rs2, off; pseudo: blez (rs1=zero), bgez (rs2=zero).
type Bge struct {
	base

	rs1, rs2 string
	off      imm // pc-relative byte offset
}

// Bge - bge rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bge(rs1, rs2 Reg, off int64) Instr {
	return Bge{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

func decodeBge(w uint32) Instr {
	return Bge{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func (i Bge) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.rs1 == "zero" {
		return fmt.Sprintf("blez %s, %s", i.rs2, target.text())
	}

	if i.rs2 == "zero" {
		return fmt.Sprintf("bgez %s, %s", i.rs1, target.text())
	}

	return fmt.Sprintf("bge %s, %s, %s", i.rs1, i.rs2, target.text())
}

func (i Bge) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bge"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func newBge(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bge")
	if err != nil {
		return nil, err
	}

	return Bge{
		rs1: rs1,
		rs2: rs2,
		off: t,
	}, nil
}
