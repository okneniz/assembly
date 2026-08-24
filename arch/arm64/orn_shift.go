package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// OrnShift - orn rd, rn, rm[, shift #imm6]; pseudo: mvn (Rn = zr).
type OrnShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

const (
	OrnShiftX uint32 = 0xAA200000
	OrnShiftW uint32 = 0x2A200000
)

func decodeOrnShift(w uint32, addr uint64) Instr {
	return OrnShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

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

func (i OrnShift) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i OrnShift) MarshalJSON() ([]byte, error) {
	return i.marshal("orn", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
