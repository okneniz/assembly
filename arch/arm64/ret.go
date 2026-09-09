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

// newRet - the Ret constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newRet(b base, rn string) Ret {
	return Ret{
		base: b,
		rn:   rn,
	}
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

// Ret — ret rn (ret without an operand is Ret with the x30 register: X(30)).
func (Builder) Ret(rn Reg) (Instr, error) {
	if err := requireClass(
		rn,
		"Ret",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return newRet(base{}, rn.name()), nil
}

func decodeRet(w uint32) Instr {
	return newRet(newBase(w), regNameX(w>>5&0x1f))
}
