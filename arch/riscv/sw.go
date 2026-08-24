package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sw — sw rs2, off(rs1).
type Sw struct {
	base

	rs1, rs2 string
	off      imm
}

// NewSw — sw rs2, off(rs1).
func NewSw(rs2, rs1 Reg, off Off) Instr {
	return Sw{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

func decodeSw(w uint32, addr uint64) Instr {
	return Sw{
		base: newBase(addr, w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

// cSw - compressed forms (c.sw/c.swsp): base - halfword, length 2.
func cSw(h uint32, addr uint64, rs1, rs2 string, off int64) Sw {
	return Sw{
		base: newHalfBase(h, addr),
		rs1:  rs1,
		rs2:  rs2,
		off:  immNum(off),
	}
}

func (i Sw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sw %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Sw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["sw"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	if !o.NoRVC {
		if rs2, ok := asmRegNum[i.rs2]; ok && !rs2.fp && i.rs1 == "sp" && immFits(off, 0, 252, 4) {
			u := uint16(off)
			// c.swsp: u[5:2]=h[12:9], u[7:6]=h[8:7]
			return writeHalf(w, 0xC002|uint16(rs2.num)<<2|(u>>2&0xf)<<9|(u>>6&3)<<7)
		}

		if r3ok(i.rs1) && r3ok(i.rs2) && immFits(off, 0, 124, 4) {
			u := uint16(off)
			// c.sw: u[5:3]=h[12:10], u[2]=h6, u[6]=h5
			return writeHalf(w, 0xC000|cr3(i.rs1)<<7|cr3(i.rs2)<<2|
				(u>>3&7)<<10|(u>>2&1)<<6|(u>>6&1)<<5)
		}
	}

	return writeWord(w, word)
}

func (i Sw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSw(ops []Op) (Instr, error) {
	rs2, base, off, err := wantR2M(ops, "sw", false)
	if err != nil {
		return nil, err
	}

	return Sw{
		rs1: base,
		rs2: rs2,
		off: off,
	}, nil
}
