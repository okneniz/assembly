package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulhDu - mulh.du rd, rj, rk (3R): rd = high half of the unsigned 128-bit product rj * rk.
type MulhDu struct {
	base

	rd, rj, rk uint8
}

// MulhDu - mulh.du rd, rj, rk.
func (Builder) MulhDu(rd, rj, rk Reg) Instr {
	return MulhDu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulhDu(w uint32) Instr {
	return MulhDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulhDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh.du %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulhDu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["mulh.du"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
