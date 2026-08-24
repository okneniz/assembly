package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ret - ret [xn] (defaults to x30).
type Ret struct {
	base

	rn string
}

const retMatch = 0xD65F0000

// NewRet - ret rn (ret without an operand is NewRet with the x30 register: X(30)).
func NewRet(rn Reg) (Instr, error) {
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

	return Ret{rn: rn.name()}, nil
}

func (i Ret) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == "x30" {
		return "ret"
	}

	return "ret " + i.rn
}

func (i Ret) Encode(w io.Writer, pc uint64) (int64, error) {
	num, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("ret: %w", err)
	}

	return writeWord(w, retMatch|num<<5)
}

func (i Ret) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ret",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Branch",
		map[string]any{"Rn": i.rn},
	)
}

func decodeRet(w uint32, addr uint64) Instr {
	return Ret{
		base: newBase(addr, w),
		rn:   regNameX(w >> 5 & 0x1f),
	}
}
