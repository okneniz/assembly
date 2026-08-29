package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulhD - mulh.d rd, rj, rk (3R): rd = high half of the signed 128-bit product rj * rk.
type MulhD struct {
	base

	rd, rj, rk uint8
}

// NewMulhD - mulh.d rd, rj, rk.
func NewMulhD(rd, rj, rk Reg) Instr {
	return MulhD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulhD(w uint32, addr uint64) Instr {
	return MulhD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulhD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulhD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mulh.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulhD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mulh.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
