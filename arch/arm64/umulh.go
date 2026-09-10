package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Umulh — umulh rd, rn, rm.
type Umulh struct {
	base

	rd, rn, rm string
}

// newUmulh - the Umulh constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newUmulh(b base, rd Reg, rn Reg, rm Reg) (Umulh, error) {
	err := requireClass(
		rd,
		"Umulh",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Umulh{}, err
	}

	err = requireClass(
		rn,
		"Umulh",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Umulh{}, err
	}

	err = requireClass(
		rm,
		"Umulh",
		"rm",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Umulh{}, err
	}

	return Umulh{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const UmulhX uint32 = 0x9BC07C00

func (i Umulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("umulh %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Umulh) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, UmulhX, 0)
	if err != nil {
		return 0, fmt.Errorf("umulh: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("umulh: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Umulh(rd, rn, rm Reg) (Instr, error) {
	return newUmulh(base{}, rd, rn, rm)
}

func decodeUmulh(w uint32) (Instr, error) {
	in, err := newUmulh(
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
