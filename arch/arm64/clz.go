package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Clz — clz rd, rn.
type Clz struct {
	base

	rd, rn string
}

// newClz - the Clz constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newClz(b base, rd Reg, rn Reg) (Clz, error) {
	err := requireClass(
		rd,
		"Clz",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Clz{}, err
	}

	err = requireClass(
		rn,
		"Clz",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Clz{}, err
	}

	err = requireWidth(
		"Clz",
		rd,
		rn,
	)

	if err != nil {
		return Clz{}, err
	}

	return Clz{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const ClzX uint32 = 0xDAC01000

func (i Clz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clz %s, %s", i.rd, i.rn)
}

func (i Clz) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, ClzX, 0x5AC01000)
	if err != nil {
		return 0, fmt.Errorf("clz: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("clz: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Clz(rd, rn Reg) (Instr, error) {
	return newClz(base{}, rd, rn)
}

func decodeClz(w uint32) (Instr, error) {
	in, err := newClz(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
