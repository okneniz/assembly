package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LsrReg — lsr rd, rn, rm.
type LsrReg struct {
	rd, rn, rm string
}

// newLsrReg - the LsrReg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLsrReg(rd Reg, rn Reg, rm Reg) (LsrReg, error) {
	err := requireClass(
		rd,
		"LsrReg",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LsrReg{}, err
	}

	err = requireClass(
		rn,
		"LsrReg",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LsrReg{}, err
	}

	err = requireClass(
		rm,
		"LsrReg",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LsrReg{}, err
	}

	err = requireWidth(
		"LsrReg",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return LsrReg{}, err
	}

	return LsrReg{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const LsrRegX uint32 = 0x9AC02400

func (i LsrReg) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, LsrRegX, 0x1AC02400)
	if err != nil {
		return 0, fmt.Errorf("lsr: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("lsr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i LsrReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lsr %s, %s, %s", i.rd, i.rn, i.rm)
}
