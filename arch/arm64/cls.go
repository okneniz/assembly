package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cls — cls rd, rn.
type Cls struct {
	base

	rd, rn string
}

// newCls - the Cls constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newCls(b base, rd Reg, rn Reg) (Cls, error) {
	err := requireClass(
		rd,
		"Cls",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Cls{}, err
	}

	err = requireClass(
		rn,
		"Cls",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Cls{}, err
	}

	err = requireWidth(
		"Cls",
		rd,
		rn,
	)

	if err != nil {
		return Cls{}, err
	}

	return Cls{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const ClsX uint32 = 0xDAC01400

func (i Cls) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cls %s, %s", i.rd, i.rn)
}

func (i Cls) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, ClsX, 0x5AC01400)
	if err != nil {
		return 0, fmt.Errorf("cls: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("cls: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Cls(rd, rn Reg) (Instr, error) {
	return newCls(base{}, rd, rn)
}

func decodeCls(w uint32) (Instr, error) {
	in, err := newCls(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
