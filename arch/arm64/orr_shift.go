package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// OrrShift — orr rd, rn, rm[, shift #imm6]; pseudo: mov (Rn = zr).
type OrrShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

const (
	OrrShiftX uint32 = 0xAA000000
	OrrShiftW uint32 = 0x2A000000
)

// OrrShift — orr rd, rn, rm[, shift #imm6] (mov when Rn = zr).
// Register 31 reads as zr (SP/WSP are not allowed — use XZR/WZR);
// shift — lsl/lsr/asr/ror.
func (Builder) OrrShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "OrrShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "OrrShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "OrrShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("OrrShift", rd, rn, rm); err != nil {
		return nil, err
	}

	return OrrShift{
		rd:    rd.name(),
		rn:    rn.name(),
		rm:    rm.name(),
		imm6:  imm.v,
		shift: sh.String(),
		isf:   rd.Is64(),
	}, nil
}

func decodeOrrShift(w uint32, addr uint64) Instr {
	return OrrShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

func (i OrrShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rn == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("mov %s, %s, %s #%d", i.rd, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("mov %s, %s", i.rd, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("orr %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("orr %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i OrrShift) Encode(w io.Writer, pc uint64) (int64, error) {
	match := OrrShiftX
	if !i.isf {
		match = OrrShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("orr: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (i OrrShift) MarshalJSON() ([]byte, error) {
	return i.marshal("orr", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
