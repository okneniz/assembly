package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddImm — add rd, rn, #imm12[, lsl #12]; pseudo: mov (imm=0, Rd≠Rn).
type AddImm struct {
	base

	rdNum, rnNum uint32 // 31: sp/wsp
	imm12        uint32
	shift        bool // lsl #12
	isf          bool
}

// newAddImm - the AddImm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAddImm(b base, rd Reg, rn Reg, imm Imm12, sh Sh12) (AddImm, error) {
	err := requireClass(
		rd,
		"AddImm",
		"rd",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddImm{}, err
	}

	err = requireClass(
		rn,
		"AddImm",
		"rn",
		"register 31 reads as sp/wsp — use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return AddImm{}, err
	}

	err = requireWidth(
		"AddImm",
		rd,
		rn,
	)

	if err != nil {
		return AddImm{}, err
	}

	return AddImm{
		base:  b,
		rdNum: rd.bits(),
		rnNum: rn.bits(),
		imm12: imm.v,
		shift: sh == LSL12,
		isf:   rd.Is64(),
	}, nil
}

const (
	AddImmX uint32 = 0x91000000
	AddImmW uint32 = 0x11000000
)

func (i AddImm) ObjDump(_ disasm.ViewCtx) string {
	rd := addSubRegName(i.rdNum, i.isf, false)
	rn := addSubRegName(i.rnNum, i.isf, false)
	imm := fmt.Sprintf("#0x%x", i.imm12)
	// LLVM alias: imm=0 and (rd≠rn or an sp/wsp pair — mov sp, sp);
	// 32-bit sp-sp stays add
	if !i.shift && i.imm12 == 0 && (i.rdNum != i.rnNum || (i.rdNum == 31 && i.isf)) {
		return fmt.Sprintf("mov %s, %s", rd, rn)
	}

	if i.shift {
		return fmt.Sprintf("add %s, %s, %s, lsl #12", rd, rn, imm)
	}

	return fmt.Sprintf("add %s, %s, %s", rd, rn, imm)
}

func (i AddImm) Encode(w io.Writer) (int64, error) {
	match := AddImmX
	if !i.isf {
		match = AddImmW
	}

	if i.imm12 > 0xfff {
		return 0, errors.New("add: imm12 out of range")
	}

	sh := uint32(0)
	if i.shift {
		sh = 1
	}

	return writeWord(w, match|i.rdNum|i.rnNum<<5|i.imm12<<10|sh<<22)
}

func (Builder) AddImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	return newAddImm(base{}, rd, rn, imm, sh)
}

func decodeAddImm(w uint32) (Instr, error) {
	in, err := newAddImm(newBase(w),
		numReg(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), imm12Of(w>>10&0xfff), sh12Of(w>>22&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
