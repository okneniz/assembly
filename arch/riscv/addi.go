package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Addi - addi rd, rs1, imm; pseudo: nop (zero,zero,0), li (rs1=zero),
// mv (imm=0).
type Addi struct {
	base

	rd, rs1 string
	imm     imm
}

// NewAddi - addi rd, rs1, imm (rs1 = zero is printed as li, imm = 0 - mv).
func NewAddi(rd, rs1 Reg, imm Imm12) Instr {
	return Addi{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

func decodeAddi(w uint32, addr uint64) Instr {
	return Addi{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

// cAddi - compressed forms (c.addi/c.nop/c.li/c.addi4spn/c.addi16sp): base - halfword, length 2.
func cAddi(h uint32, addr uint64, rd, rs1 string, imm int64) Addi {
	return Addi{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		imm:  immNum(imm),
	}
}

func (i Addi) ObjDump(_ disasm.ViewCtx) string {
	switch {
	case i.rd == "zero" && i.rs1 == "zero" && i.imm.val == 0:
		return "nop"
	case i.rs1 == "zero":
		return fmt.Sprintf("li %s, %s", i.rd, i.imm.text())
	case i.imm.val == 0:
		return fmt.Sprintf("mv %s, %s", i.rd, i.rs1)
	}

	return fmt.Sprintf("addi %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Addi) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["addi"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if !o.NoRVC {
		if half, ok := compressAddi(i.rd, i.rs1, v); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

func (i Addi) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"addi",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

// compressAddi - c.nop/c.li/c.addi4spn/c.addi/c.addi16sp (c.addi
// is preferred over c.addi16sp, as in GAS).
func compressAddi(rd, rs1 string, imm int64) (uint16, bool) {
	r := r5(rd)
	switch {
	case rd == "zero" && rs1 == "zero" && imm == 0:
		return 0x0001, true // c.nop
	case rs1 == "zero" && r != 0 && fits6(imm):
		return 0x4001 | r<<7 | ciBits(imm), true // c.li
	case rs1 == "sp" && r3ok(rd) && immFits(imm, 4, 1020, 4):
		// c.addi4spn: u[9:6]=h[10:7], u[5:4]=h[12:11], u[3]=h[5], u[2]=h[6]
		u := uint16(imm)
		return 0x0000 | cr3(rd)<<2 | (u>>6&0xf)<<7 | (u>>4&3)<<11 | (u>>3&1)<<5 | (u>>2&1)<<6, true
	case rd == rs1 && r != 0 && fits6(imm):
		return 0x0001 | r<<7 | ciBits(imm), true // c.addi - any rd (including sp)
	case rd == "sp" && rs1 == "sp" && imm != 0 && imm%16 == 0 && imm >= -496 && imm <= 496:
		// c.addi16sp: u[9]=h12, u[8:7]=h[4:3], u[6]=h5, u[5]=h2, u[4]=h6
		u := uint16(imm & 0x3ff)
		h := uint16(0x6101) // Q1 funct3=011, rd=sp
		h |= (u>>9&1)<<12 | (u>>7&3)<<3 | (u>>6&1)<<5 | (u>>5&1)<<2 | (u>>4&1)<<6
		return h, true
	}

	return 0, false
}

func newAddi(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "addi")
	if err != nil {
		return nil, err
	}

	return Addi{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
