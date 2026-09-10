package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RorReg — ror rd, rn, rm.
type RorReg struct {
	base

	rd, rn, rm string
}

// newRorReg - the RorReg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newRorReg(b base, rd Reg, rn Reg, rm Reg) (RorReg, error) {
	err := requireClass(
		rd,
		"RorReg",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return RorReg{}, err
	}

	err = requireClass(
		rn,
		"RorReg",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return RorReg{}, err
	}

	err = requireClass(
		rm,
		"RorReg",
		"rm",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return RorReg{}, err
	}

	return RorReg{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const RorRegX uint32 = 0x9A002C00

func (i RorReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ror %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i RorReg) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, RorRegX, 0)
	if err != nil {
		return 0, fmt.Errorf("ror: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ror: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) RorReg(rd, rn, rm Reg) (Instr, error) {
	return newRorReg(base{}, rd, rn, rm)
}

func decodeRorReg(w uint32) (Instr, error) {
	in, err := newRorReg(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
