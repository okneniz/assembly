package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubImm — sub rd, rn, #imm12[, lsl #12].
type SubImm struct {
	base

	rdNum, rnNum uint32 // 31: sp/wsp
	imm12        uint32
	shift        bool // lsl #12
	isf          bool
}

// newSubImm - the SubImm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSubImm(b base, rd Reg, rn Reg, imm Imm12, sh Sh12) (SubImm, error) {
	err := requireClass(
		rd,
		"SubImm",
		"rd",
		"register 31 reads as sp/wsp - use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubImm{}, err
	}

	err = requireClass(
		rn,
		"SubImm",
		"rn",
		"register 31 reads as sp/wsp - use SP/WSP",
		classX,
		classW,
		classSP,
		classWSP,
	)

	if err != nil {
		return SubImm{}, err
	}

	err = requireWidth(
		"SubImm",
		rd,
		rn,
	)

	if err != nil {
		return SubImm{}, err
	}

	return SubImm{
		base:  b,
		rdNum: rd.bits(),
		rnNum: rn.bits(),
		imm12: imm.v,
		shift: sh == LSL12,
		isf:   rd.Is64(),
	}, nil
}

const (
	SubImmX uint32 = 0xD1000000
	SubImmW uint32 = 0x51000000
)

func (i SubImm) ObjDump(_ disasm.ViewCtx) string {
	rd := addSubRegName(i.rdNum, i.isf, false)
	rn := addSubRegName(i.rnNum, i.isf, false)
	imm := fmt.Sprintf("#0x%x", i.imm12)

	if i.shift {
		return fmt.Sprintf("sub %s, %s, %s, lsl #12", rd, rn, imm)
	}

	return fmt.Sprintf("sub %s, %s, %s", rd, rn, imm)
}

func (i SubImm) Encode(w io.Writer) (int64, error) {
	match := SubImmX
	if !i.isf {
		match = SubImmW
	}

	if i.imm12 > 0xfff {
		return 0, errors.New("sub: imm12 out of range")
	}

	sh := uint32(0)
	if i.shift {
		sh = 1
	}

	return writeWord(w, match|i.rdNum|i.rnNum<<5|i.imm12<<10|sh<<22)
}

func (Builder) SubImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	return newSubImm(base{}, rd, rn, imm, sh)
}

func decodeSubImm(w uint32) (Instr, error) {
	in, err := newSubImm(newBase(w),
		numReg(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), imm12Of(w>>10&0xfff), sh12Of(w>>22&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
