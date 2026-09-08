package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lui - lui rd, imm.
type Lui struct {
	base

	rd  string
	imm imm
}

// Lui - lui rd, imm20.
func (Builder) Lui(rd Reg, imm Imm20) Instr {
	return Lui{
		rd:  rd.name(),
		imm: immNum(imm.v),
	}
}

func decodeLui(w uint32, addr uint64) Instr {
	return Lui{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		imm:  immNum(int64(uImm(w))),
	}
}

// cLui - compressed forms (c.lui): base - halfword, length 2.
func cLui(h uint32, addr uint64, rd string, imm int64) Lui {
	return Lui{
		base: newHalfBase(h, addr),
		rd:   rd,
		imm:  immNum(imm),
	}
}

func (i Lui) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lui %s, %s", i.rd, i.imm.text())
}

func (i Lui) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encU(v)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lui"][0] | regBits(i.rd)<<7 | bits
	if !o.NoRVC {
		if half, ok := compressLui(i.rd, v); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

// compressLui — c.lui: nzimm6 ≠ 0, rd ∉ {x0, sp}. The immediate is the
// sign-extended low 20 bits of v (lui's field): the 0xffff8 and -8
// spellings of the same value compress identically (negative pages,
// e.g. li rd, -0x8000 → c.lui rd, -8).
func compressLui(rd string, v int64) (uint16, bool) {
	r := r5(rd)
	sv := int64(int32(uint32(v)<<12)) >> 12 // sext(v[19:0])
	if r == 0 || r == 2 || sv == 0 || !fits6(sv) {
		return 0, false
	}

	return 0x6001 | r<<7 | ciBits(sv), true
}

// newLui - constructor from parsing: lui rd, imm.
func newLui(ops []Op) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("lui: want rd, imm")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("lui: %w", err)
	}

	e, err := wantExpr(ops[1])
	if err != nil {
		return nil, fmt.Errorf("lui: %w", err)
	}

	return Lui{
		rd:  rd,
		imm: e,
	}, nil
}
