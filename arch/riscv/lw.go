package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lw - lw rd, off(rs1).
type Lw struct {
	base

	rd, rs1 string
	off     imm
}

// Lw - lw rd, off(rs1).
func (Builder) Lw(rd, rs1 Reg, off Off) Instr {
	return Lw{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeLw(w uint32, addr uint64) Instr {
	return Lw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

// cLw - compressed forms (c.lw/c.lwsp): base - halfword, length 2.
func cLw(h uint32, addr uint64, rd, rs1 string, off int64) Lw {
	return Lw{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		off:  immNum(off),
	}
}

func (i Lw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lw %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Lw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if !o.NoRVC {
		if r := r5(i.rd); i.rs1 == "sp" && r != 0 && immFits(off, 0, 252, 4) {
			u := uint16(off)
			// c.lwsp: u[5]=h12, u[4:2]=h[6:4], u[7:6]=h[3:2]
			return writeHalf(w, 0x4002|r<<7|(u>>5&1)<<12|(u>>2&7)<<4|(u>>6&3)<<2)
		}

		if r3ok(i.rd) && r3ok(i.rs1) && immFits(off, 0, 124, 4) {
			u := uint16(off)
			// c.lw: u[5:3]=h[12:10], u[2]=h6, u[6]=h5
			return writeHalf(w, 0x4000|cr3(i.rs1)<<7|cr3(i.rd)<<2|
				(u>>3&7)<<10|(u>>2&1)<<6|(u>>6&1)<<5)
		}
	}

	return writeWord(w, word)
}

func (i Lw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newLw(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lw", false)
	if err != nil {
		return nil, err
	}

	return Lw{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
