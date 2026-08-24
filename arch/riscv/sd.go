package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sd — sd rs2, off(rs1).
type Sd struct {
	base

	rs1, rs2 string
	off      imm
}

// NewSd — sd rs2, off(rs1).
func NewSd(rs2, rs1 Reg, off Off) Instr {
	return Sd{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

func decodeSd(w uint32, addr uint64) Instr {
	return Sd{
		base: newBase(addr, w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

// cSd - compressed forms (c.sd/c.sdsp): base - halfword, length 2.
func cSd(h uint32, addr uint64, rs1, rs2 string, off int64) Sd {
	return Sd{
		base: newHalfBase(h, addr),
		rs1:  rs1,
		rs2:  rs2,
		off:  immNum(off),
	}
}

func (i Sd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sd %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Sd) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["sd"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	if !o.NoRVC {
		if rs2, ok := asmRegNum[i.rs2]; ok && !rs2.fp && i.rs1 == "sp" && immFits(off, 0, 56, 8) {
			u := uint16(off)
			// c.sdsp: u[5:3]=h[12:10]
			return writeHalf(w, 0xE002|uint16(rs2.num)<<2|(u>>3&7)<<10)
		}

		if r3ok(i.rs1) && r3ok(i.rs2) && immFits(off, 0, 248, 8) {
			u := uint16(off)
			// c.sd: u[5:3]=h[12:10], u[7:6]=h[6:5]
			return writeHalf(w, 0xE000|cr3(i.rs1)<<7|cr3(i.rs2)<<2|
				(u>>3&7)<<10|(u>>6&3)<<5)
		}
	}

	return writeWord(w, word)
}

func (i Sd) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sd",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSd(ops []Op) (Instr, error) {
	rs2, base, off, err := wantR2M(ops, "sd", false)
	if err != nil {
		return nil, err
	}

	return Sd{
		rs1: base,
		rs2: rs2,
		off: off,
	}, nil
}
