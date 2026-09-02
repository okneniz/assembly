package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubsImm — subs rd, rn, #imm12[, lsl #12]; pseudo: cmp (Rd = zr).
type SubsImm struct {
	base

	rdNum, rnNum uint32 // 31: sp/wsp (Rd with S=1 - zr)
	imm12        uint32
	shift        bool // lsl #12
	isf          bool
}

const (
	SubsImmX uint32 = 0xF1000000
	SubsImmW uint32 = 0x71000000
)

// SubsImm — subs rd, rn, #imm12[, lsl #12] (cmp when Rd = zr).
// Rd: register 31 reads as zr; Rn — as sp/wsp.
func (Builder) SubsImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	if err := requireClass(
		rd,
		"SubsImm",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the cmp form)",
		classX,
		classW,
		classXZR,
		classWZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "SubsImm", "rn", "register 31 reads as sp/wsp — use SP/WSP",
		classX, classW, classSP, classWSP); err != nil {
		return nil, err
	}

	if err := requireWidth("SubsImm", rd, rn); err != nil {
		return nil, err
	}

	return SubsImm{
		rdNum: rd.bits(),
		rnNum: rn.bits(),
		imm12: imm.v,
		shift: sh == LSL12,
		isf:   rd.Is64(),
	}, nil
}

func decodeSubsImm(w uint32, addr uint64) Instr {
	return SubsImm{
		base:  newBase(addr, w),
		rdNum: w & 0x1f,
		rnNum: w >> 5 & 0x1f,
		imm12: w >> 10 & 0xfff,
		shift: w>>22&1 == 1,
		isf:   w>>31&1 == 1,
	}
}

func (i SubsImm) ObjDump(_ disasm.ViewCtx) string {
	rd := addSubRegName(i.rdNum, i.isf, true)
	rn := addSubRegName(i.rnNum, i.isf, false)
	imm := fmt.Sprintf("#0x%x", i.imm12)
	if i.rdNum == 31 {
		rnz := addSubRegName(i.rnNum, i.isf, true)
		if i.shift {
			return fmt.Sprintf("cmp %s, %s, lsl #12", rnz, imm)
		}

		return fmt.Sprintf("cmp %s, %s", rnz, imm)
	}

	if i.shift {
		return fmt.Sprintf("subs %s, %s, %s, lsl #12", rd, rn, imm)
	}

	return fmt.Sprintf("subs %s, %s, %s", rd, rn, imm)
}

func (i SubsImm) Encode(w io.Writer, pc uint64) (int64, error) {
	match := SubsImmX
	if !i.isf {
		match = SubsImmW
	}

	if i.imm12 > 0xfff {
		return 0, errors.New("subs: imm12 out of range")
	}

	sh := uint32(0)
	if i.shift {
		sh = 1
	}

	return writeWord(w, match|i.rdNum|i.rnNum<<5|i.imm12<<10|sh<<22)
}

func (i SubsImm) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"subs",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing - immediate",
		map[string]any{
			"Rd":    addSubRegName(i.rdNum, i.isf, true),
			"Rn":    addSubRegName(i.rnNum, i.isf, false),
			"imm12": i.imm12,
		},
	)
}
