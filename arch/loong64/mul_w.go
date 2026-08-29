package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulW - mul.w rd, rj, rk (3R): rd = sign32(low32(rj) * low32(rk)).
type MulW struct {
	base

	rd, rj, rk uint8
}

// NewMulW - mul.w rd, rj, rk.
func NewMulW(rd, rj, rk Reg) Instr {
	return MulW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulW(w uint32, addr uint64) Instr {
	return MulW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mul.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mul.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mul.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
