package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Andi - andi rd, rs1, imm; pseudo: zext.b (imm=0xff); compression: c.andi.
type Andi struct {
	base

	rd, rs1 string
	imm     imm
}

// Andi - andi rd, rs1, imm (imm = 0xff is printed as zext.b).
func (Builder) Andi(rd, rs1 Reg, imm Imm12) Instr {
	return Andi{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

func decodeAndi(w uint32, addr uint64) Instr {
	return Andi{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

// cAndi - compressed forms (c.andi): base - halfword, length 2.
func cAndi(h uint32, addr uint64, rd, rs1 string, imm int64) Andi {
	return Andi{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		imm:  immNum(imm),
	}
}

func (i Andi) ObjDump(_ disasm.ViewCtx) string {
	if i.imm.val == 0xff {
		return fmt.Sprintf("zext.b %s, %s", i.rd, i.rs1)
	}

	return fmt.Sprintf("andi %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Andi) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["andi"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if !o.NoRVC {
		// decoder mirror: c.andi - bits[11:10]=01 (see decodeCA)
		if r3ok(i.rd) && i.rd == i.rs1 && fits6(v) {
			return writeHalf(w, 0x8401|cr3(i.rs1)<<7|ciBits(v))
		}
	}

	return writeWord(w, word)
}

func (i Andi) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"andi",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newAndi(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "andi")
	if err != nil {
		return nil, err
	}

	return Andi{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
