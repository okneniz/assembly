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

// newSubImm - the SubImm constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newSubImm(b base, rdNum uint32, rnNum uint32, imm12 uint32, shift bool, isf bool) SubImm {
	return SubImm{
		base:  b,
		rdNum: rdNum,
		rnNum: rnNum,
		imm12: imm12,
		shift: shift,
		isf:   isf,
	}
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

// SubImm — sub rd, rn, #imm12[, lsl #12]. Register 31 reads as sp/wsp
// (XZR/WZR are not allowed - use SP/WSP).
func (Builder) SubImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	if err := requireClass(rd, "SubImm", "rd", "register 31 reads as sp/wsp - use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "SubImm", "rn", "register 31 reads as sp/wsp - use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("SubImm", rd, rn); err != nil {
		return nil, err
	}

	return newSubImm(base{}, rd.bits(), rn.bits(), imm.v, sh == LSL12, rd.Is64()), nil
}

func decodeSubImm(w uint32) Instr {
	return newSubImm(newBase(w), w&0x1f, w>>5&0x1f, w>>10&0xfff, w>>22&1 == 1, w>>31&1 == 1)
}
