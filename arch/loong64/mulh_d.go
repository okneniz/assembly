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

// MulhD - mulh.d rd, rj, rk.
func (Builder) MulhD(rd, rj, rk Reg) Instr {
	return MulhD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulhD(w uint32) Instr {
	return MulhD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulhD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulhD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mulh.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
