package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Br — br xn (indirect branch).
type Br struct {
	base

	rn string
}

// newBr - the Br constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newBr(b base, rn Reg) (Br, error) {
	err := requireClass(
		rn,
		"Br",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Br{}, err
	}

	return Br{
		base: b,
		rn:   rn.name(),
	}, nil
}

const brMatch = 0xD61F0000

func (i Br) ObjDump(_ disasm.ViewCtx) string {
	return "br " + i.rn
}

func (i Br) Encode(w io.Writer) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("br: %w", err)
	}

	return writeWord(w, brMatch|num<<5)
}

func (Builder) Br(rn Reg) (Instr, error) {
	return newBr(base{}, rn)
}

func decodeBr(w uint32) (Instr, error) {
	in, err := newBr(newBase(w), gprOf(w>>5&0x1f, true))
	if err != nil {
		return nil, err
	}

	return in, nil
}
