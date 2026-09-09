package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Blr — blr xn (indirect call).
type Blr struct {
	base

	rn string
}

// newBlr - the Blr constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newBlr(b base, rn string) Blr {
	return Blr{
		base: b,
		rn:   rn,
	}
}

const blrMatch = 0xD63F0000

func (i Blr) ObjDump(_ disasm.ViewCtx) string {
	return "blr " + i.rn
}

func (i Blr) Encode(w io.Writer) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("blr: %w", err)
	}

	return writeWord(w, blrMatch|num<<5)
}

// Blr — blr rn (indirect call). Only x registers (X/XZR).
func (Builder) Blr(rn Reg) (Instr, error) {
	if err := requireClass(
		rn,
		"Blr",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return newBlr(base{}, rn.name()), nil
}

func decodeBlr(w uint32) Instr {
	return newBlr(newBase(w), regNameX(w>>5&0x1f))
}
