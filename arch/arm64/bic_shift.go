package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BicShift — bic rd, rn, rm[, shift #imm6].
type BicShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

const (
	BicShiftX uint32 = 0x8A200000
	BicShiftW uint32 = 0x0A200000
)

// BicShift — bic rd, rn, rm[, shift #imm6]. Register 31 reads as
// zr (SP/WSP are not allowed — use XZR/WZR); shift — lsl/lsr/asr/ror.
func (Builder) BicShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "BicShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "BicShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "BicShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("BicShift", rd, rn, rm); err != nil {
		return nil, err
	}

	return BicShift{
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

func decodeBicShift(w uint32, addr uint64) Instr {
	return BicShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

func (i BicShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("bic %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("bic %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i BicShift) Encode(w io.Writer, pc uint64) (int64, error) {
	match := BicShiftX
	if !i.isf {
		match = BicShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bic: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("bic: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("bic: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (i BicShift) MarshalJSON() ([]byte, error) {
	return i.marshal("bic", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
