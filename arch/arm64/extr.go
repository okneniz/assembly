package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Extr — extr rd, rn, rm, #lsb; pseudo: ror rd, rn, #imm (rn == rm).
type Extr struct {
	base

	rd, rn, rm string
	lsb        uint32
}

// newExtr - the Extr constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newExtr(b base, rd string, rn string, rm string, lsb uint32) Extr {
	return Extr{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
		lsb:  lsb,
	}
}

const extrX uint32 = 0x93000000

func (i Extr) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == i.rm {
		return fmt.Sprintf("ror %s, %s, #0x%x", i.rd, i.rn, i.lsb)
	}

	return fmt.Sprintf("extr %s, %s, %s, #0x%x", i.rd, i.rn, i.rm, i.lsb)
}

func (i Extr) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("extr: %w", err)
	}

	if i.lsb > 63 {
		return 0, fmt.Errorf("extr: lsb %#x out of range", i.lsb)
	}

	return writeWord(w, extrX|rd|rn<<5|i.lsb<<10|rm<<16)
}

// Extr — extr rd, rn, rm, #lsb (ror when Rn == Rm). Only the
// 64-bit form; register 31 reads as zr (SP/WSP are not allowed — use XZR);
// lsb — 0..63.
func (Builder) Extr(rd, rn, rm Reg, lsb Imm6) (Instr, error) {
	if err := requireClass(rd, "Extr", "rd",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Extr", "rn",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "Extr", "rm",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	return newExtr(base{}, rd.name(), rn.name(), rm.name(), lsb.v), nil
}

func decodeExtr(w uint32) Instr {
	return newExtr(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
		w>>10&0x3f,
	)
}
