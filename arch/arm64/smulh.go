package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Smulh — smulh rd, rn, rm.
type Smulh struct {
	base

	rd, rn, rm string
}

// newSmulh - the Smulh constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSmulh(b base, rd Reg, rn Reg, rm Reg) (Smulh, error) {
	err := requireClass(
		rd,
		"Smulh",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Smulh{}, err
	}

	err = requireClass(
		rn,
		"Smulh",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Smulh{}, err
	}

	err = requireClass(
		rm,
		"Smulh",
		"rm",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Smulh{}, err
	}

	return Smulh{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const SmulhX uint32 = 0x9B407C00

func (i Smulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("smulh %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Smulh) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, SmulhX, 0)
	if err != nil {
		return 0, fmt.Errorf("smulh: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("smulh: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Smulh(rd, rn, rm Reg) (Instr, error) {
	return newSmulh(base{}, rd, rn, rm)
}

func decodeSmulh(w uint32) (Instr, error) {
	in, err := newSmulh(
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
