package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Udiv — udiv rd, rn, rm.
type Udiv struct {
	rd, rn, rm string
}

// newUdiv - the Udiv constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newUdiv(rd Reg, rn Reg, rm Reg) (Udiv, error) {
	err := requireClass(
		rd,
		"Udiv",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Udiv{}, err
	}

	err = requireClass(
		rn,
		"Udiv",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Udiv{}, err
	}

	err = requireClass(
		rm,
		"Udiv",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Udiv{}, err
	}

	err = requireWidth(
		"Udiv",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return Udiv{}, err
	}

	return Udiv{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const UdivX uint32 = 0x9AC00800

func (i Udiv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("udiv %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Udiv) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, UdivX, 0x1AC00800)
	if err != nil {
		return 0, fmt.Errorf("udiv: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("udiv: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}
