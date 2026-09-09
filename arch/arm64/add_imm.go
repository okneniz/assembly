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

// newAddImm - the AddImm constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAddImm(b base, rdNum uint32, rnNum uint32, imm12 uint32, shift bool, isf bool) AddImm {
	return AddImm{
		base:  b,
		rdNum: rdNum,
		rnNum: rnNum,
		imm12: imm12,
		shift: shift,
		isf:   isf,
	}
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

// AddImm — add rd, rn, #imm12[, lsl #12]. Register 31 reads as
// sp/wsp (XZR/WZR are not allowed — use SP/WSP).
func (Builder) AddImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	if err := requireClass(rd, "AddImm", "rd", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AddImm", "rn", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("AddImm", rd, rn); err != nil {
		return nil, err
	}

	return newAddImm(base{}, rd.bits(), rn.bits(), imm.v, sh == LSL12, rd.Is64()), nil
}

func decodeAddImm(w uint32) Instr {
	return newAddImm(newBase(w), w&0x1f, w>>5&0x1f, w>>10&0xfff, w>>22&1 == 1, w>>31&1 == 1)
}
