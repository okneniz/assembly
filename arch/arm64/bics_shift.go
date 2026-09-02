package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BicsShift — bics rd, rn, rm[, shift #imm6].
type BicsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

const (
	BicsShiftX uint32 = 0xEA200000
	BicsShiftW uint32 = 0x6A200000
)

// BicsShift — bics rd, rn, rm[, shift #imm6]. Register 31 reads
// as zr (SP/WSP are not allowed — use XZR/WZR); shift — lsl/lsr/asr/ror.
func (Builder) BicsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "BicsShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "BicsShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "BicsShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("BicsShift", rd, rn, rm); err != nil {
		return nil, err
	}

	return BicsShift{
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

func decodeBicsShift(w uint32, addr uint64) Instr {
	return BicsShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

func (i BicsShift) ObjDump(_ disasm.ViewCtx) string {
	if i.imm6 != 0 {
		return fmt.Sprintf("bics %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("bics %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i BicsShift) Encode(w io.Writer, pc uint64) (int64, error) {
	match := BicsShiftX
	if !i.isf {
		match = BicsShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bics: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("bics: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("bics: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (i BicsShift) MarshalJSON() ([]byte, error) {
	return i.marshal("bics", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
