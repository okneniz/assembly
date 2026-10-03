package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Blr — blr xn (indirect call).
type Blr struct {
	rn string
}

// newBlr - the Blr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newBlr(rn Reg) (Blr, error) {
	err := requireClass(
		rn,
		"Blr",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Blr{}, err
	}

	return Blr{
		rn: rn.name(),
	}, nil
}

const blrMatch = 0xD63F0000

func (i Blr) Encode(w io.Writer) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("blr: %w", err)
	}

	return writeWord(w, blrMatch|num<<5)
}

func (i Blr) ObjDump(_ disasm.ViewCtx) string {
	return "blr " + i.rn
}
