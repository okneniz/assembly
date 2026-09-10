package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ret — ret [xn] (defaults to x30).
type Ret struct {
	base

	rn string
}

// newRet - the Ret constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newRet(b base, rn Reg) (Ret, error) {
	err := requireClass(
		rn,
		"Ret",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ret{}, err
	}

	return Ret{
		base: b,
		rn:   rn.name(),
	}, nil
}

const retMatch = 0xD65F0000

func (i Ret) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == "x30" {
		return "ret"
	}

	return "ret " + i.rn
}

func (i Ret) Encode(w io.Writer) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("ret: %w", err)
	}

	return writeWord(w, retMatch|num<<5)
}

func (Builder) Ret(rn Reg) (Instr, error) {
	return newRet(base{}, rn)
}

func decodeRet(w uint32) (Instr, error) {
	in, err := newRet(newBase(w), gprOf(w>>5&0x1f, true))
	if err != nil {
		return nil, err
	}

	return in, nil
}
