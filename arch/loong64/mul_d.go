package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulD - mul.d rd, rj, rk (3R): rd = the low 64 bits of the product rj * rk.
type MulD struct {
	base

	rd, rj, rk uint8
}

// NewMulD - mul.d rd, rj, rk.
func NewMulD(rd, rj, rk Reg) Instr {
	return MulD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulD(w uint32, addr uint64) Instr {
	return MulD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mul.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mul.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mul.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
