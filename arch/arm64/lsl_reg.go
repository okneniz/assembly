package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LslReg — lsl rd, rn, rm.
type LslReg struct {
	rd, rn, rm string
}

// newLslReg - the LslReg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newLslReg(rd Reg, rn Reg, rm Reg) (LslReg, error) {
	err := requireClass(
		rd,
		"LslReg",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LslReg{}, err
	}

	err = requireClass(
		rn,
		"LslReg",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LslReg{}, err
	}

	err = requireClass(
		rm,
		"LslReg",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return LslReg{}, err
	}

	err = requireWidth(
		"LslReg",
		rd,
		rn,
		rm,
	)

	if err != nil {
		return LslReg{}, err
	}

	return LslReg{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const LslRegX uint32 = 0x9AC02000

func (i LslReg) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, LslRegX, 0x1AC02000)
	if err != nil {
		return 0, fmt.Errorf("lsl: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("lsl: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i LslReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lsl %s, %s, %s", i.rd, i.rn, i.rm)
}
