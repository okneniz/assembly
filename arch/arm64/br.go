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

// newBr - the Br constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newBr(b base, rn string) Br {
	return Br{
		base: b,
		rn:   rn,
	}
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

// Br — br xn (indirect branch). Only x registers (X/XZR).
func (Builder) Br(rn Reg) (Instr, error) {
	if err := requireClass(
		rn,
		"Br",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return newBr(base{}, rn.name()), nil
}

func decodeBr(w uint32) Instr {
	return newBr(newBase(w), regNameX(w>>5&0x1f))
}
