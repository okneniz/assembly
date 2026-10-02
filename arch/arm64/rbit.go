package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rbit — rbit rd, rn.
type Rbit struct {
	rd, rn string
}

// newRbit - the Rbit constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newRbit(rd Reg, rn Reg) (Rbit, error) {
	err := requireClass(
		rd,
		"Rbit",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Rbit{}, err
	}

	err = requireClass(
		rn,
		"Rbit",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Rbit{}, err
	}

	err = requireWidth(
		"Rbit",
		rd,
		rn,
	)

	if err != nil {
		return Rbit{}, err
	}

	return Rbit{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

const RbitX uint32 = 0xDAC00000

func (i Rbit) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rbit %s, %s", i.rd, i.rn)
}

func (i Rbit) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, RbitX, 0x5AC00000)
	if err != nil {
		return 0, fmt.Errorf("rbit: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rbit: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}
