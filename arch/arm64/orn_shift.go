package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// OrnShift — orn rd, rn, rm[, shift #imm6]; pseudo: mvn (Rn = zr).
type OrnShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

// newOrnShift - the OrnShift constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newOrnShift(
	b base,
	rd string,
	rn string,
	rm string,
	imm6 uint32,
	shift string,
	isf bool,
) OrnShift {
	return OrnShift{
		base:  b,
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
}

const (
	OrnShiftX uint32 = 0xAA200000
	OrnShiftW uint32 = 0x2A200000
)

func (i OrnShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rn == zr {
		if i.imm6 != 0 {
			return fmt.Sprintf("mvn %s, %s, %s #%d", i.rd, i.rm, i.shift, i.imm6)
		}

		return fmt.Sprintf("mvn %s, %s", i.rd, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("orn %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("orn %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i OrnShift) Encode(w io.Writer) (int64, error) {
	match := OrnShiftX
	if !i.isf {
		match = OrnShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("orn: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("orn: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("orn: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

// OrnShift — orn rd, rn, rm[, shift #imm6] (mvn when Rn = zr).
// Register 31 reads as zr (SP/WSP are not allowed — use XZR/WZR);
// shift — lsl/lsr/asr/ror.
func (Builder) OrnShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	if err := requireClass(rd, "OrnShift", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "OrnShift", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "OrnShift", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("OrnShift", rd, rn, rm); err != nil {
		return nil, err
	}

	return newOrnShift(base{}, rd.name(), rn.name(), rm.name(), imm.v, sh.String(), rd.Is64()), nil
}

func decodeOrnShift(w uint32) Instr {
	return newOrnShift(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		w>>10&0x3f,
		shiftNames[w>>22&3],
		w>>31&1 == 1,
	)
}
