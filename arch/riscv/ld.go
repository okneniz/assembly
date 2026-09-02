package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ld - ld rd, off(rs1).
type Ld struct {
	base

	rd, rs1 string
	off     imm
}

// Ld - ld rd, off(rs1).
func (Builder) Ld(rd, rs1 Reg, off Off) Instr {
	return Ld{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeLd(w uint32, addr uint64) Instr {
	return Ld{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

// cLd - compressed forms (c.ld/c.ldsp): base - halfword, length 2.
func cLd(h uint32, addr uint64, rd, rs1 string, off int64) Ld {
	return Ld{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		off:  immNum(off),
	}
}

func (i Ld) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Ld) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["ld"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if !o.NoRVC {
		if r := r5(i.rd); i.rs1 == "sp" && r != 0 && immFits(off, 0, 56, 8) {
			u := uint16(off)
			// c.ldsp: u[5]=h12, u[4:3]=h[6:5]
			return writeHalf(w, 0x6002|r<<7|(u>>5&1)<<12|(u>>3&3)<<5)
		}

		if r3ok(i.rd) && r3ok(i.rs1) && immFits(off, 0, 248, 8) {
			u := uint16(off)
			// c.ld: u[5:3]=h[12:10], u[7:6]=h[6:5]
			return writeHalf(w, 0x6000|cr3(i.rs1)<<7|cr3(i.rd)<<2|
				(u>>3&7)<<10|(u>>6&3)<<5)
		}
	}

	return writeWord(w, word)
}

func (i Ld) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ld",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newLd(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "ld", false)
	if err != nil {
		return nil, err
	}

	return Ld{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
