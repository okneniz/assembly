package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Addiw - addiw rd, rs1, imm; pseudo: sext.w (imm=0); compression: c.addiw.
type Addiw struct {
	base

	rd, rs1 string
	imm     imm
}

// Addiw - addiw rd, rs1, imm (imm = 0 is printed as sext.w).
func (Builder) Addiw(rd, rs1 Reg, imm Imm12) Instr {
	return Addiw{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

func decodeAddiw(w uint32, addr uint64) Instr {
	return Addiw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

// cAddiw - compressed forms (c.addiw): base - halfword, length 2.
func cAddiw(h uint32, addr uint64, rd, rs1 string, imm int64) Addiw {
	return Addiw{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		imm:  immNum(imm),
	}
}

func (i Addiw) ObjDump(_ disasm.ViewCtx) string {
	if i.imm.val == 0 {
		return fmt.Sprintf("sext.w %s, %s", i.rd, i.rs1)
	}

	return fmt.Sprintf("addiw %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Addiw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["addiw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if !o.NoRVC {
		if r := r5(i.rd); i.rd == i.rs1 && r != 0 && fits6(v) {
			return writeHalf(w, 0x2001|r<<7|ciBits(v)) // c.addiw
		}
	}

	return writeWord(w, word)
}

func (i Addiw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"addiw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newAddiw(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "addiw")
	if err != nil {
		return nil, err
	}

	return Addiw{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
