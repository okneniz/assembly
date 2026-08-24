package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AndsShift — ands rd, rn, rm[, shift #imm6]; pseudo: tst (Rd = zr, imm6 = 0).
type AndsShift struct {
	base

	rd, rn, rm string
	imm6       uint32
	shift      string
	isf        bool
}

const (
	AndsShiftX uint32 = 0xEA000000
	AndsShiftW uint32 = 0x6A000000
)

func decodeAndsShift(w uint32, addr uint64) Instr {
	return AndsShift{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:    armRegName(w>>16&0x1f, w>>31&1 == 1),
		imm6:  w >> 10 & 0x3f,
		shift: shiftNames[w>>22&3],
		isf:   w>>31&1 == 1,
	}
}

func (i AndsShift) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	if i.rd == zr && i.imm6 == 0 {
		return fmt.Sprintf("tst %s, %s", i.rn, i.rm)
	}

	if i.imm6 != 0 {
		return fmt.Sprintf("ands %s, %s, %s, %s #%d", i.rd, i.rn, i.rm, i.shift, i.imm6)
	}

	return fmt.Sprintf("ands %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AndsShift) Encode(w io.Writer, pc uint64) (int64, error) {
	match := AndsShiftX
	if !i.isf {
		match = AndsShiftW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ands: %w", err)
	}

	sh, err := shiftNumByName(i.shift)
	if err != nil {
		return 0, fmt.Errorf("ands: %w", err)
	}

	if i.imm6 > 63 {
		return 0, errors.New("ands: imm6 out of range")
	}

	return writeWord(w, match|rd|rn<<5|i.imm6<<10|rm<<16|sh<<22)
}

func (i AndsShift) MarshalJSON() ([]byte, error) {
	return i.marshal("ands", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - register",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imm6": i.imm6})
}
