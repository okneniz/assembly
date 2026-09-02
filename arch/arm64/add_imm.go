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

const (
	AddImmX uint32 = 0x91000000
	AddImmW uint32 = 0x11000000
)

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

	return AddImm{
		rdNum: rd.bits(),
		rnNum: rn.bits(),
		imm12: imm.v,
		shift: sh == LSL12,
		isf:   rd.Is64(),
	}, nil
}

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

func (i AddImm) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i AddImm) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"add",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing - immediate",
		map[string]any{
			"Rd":    addSubRegName(i.rdNum, i.isf, false),
			"Rn":    addSubRegName(i.rnNum, i.isf, false),
			"imm12": i.imm12,
		},
	)
}

func decodeAddImm(w uint32, addr uint64) Instr {
	return AddImm{
		base:  newBase(addr, w),
		rdNum: w & 0x1f,
		rnNum: w >> 5 & 0x1f,
		imm12: w >> 10 & 0xfff,
		shift: w>>22&1 == 1,
		isf:   w>>31&1 == 1,
	}
}
