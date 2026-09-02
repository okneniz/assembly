package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulwDW - mulw.d.w rd, rj, rk (3R): rd = sign64(low32(rj) * low32(rk)).
type MulwDW struct {
	base

	rd, rj, rk uint8
}

// MulwDW - mulw.d.w rd, rj, rk.
func (Builder) MulwDW(rd, rj, rk Reg) Instr {
	return MulwDW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulwDW(w uint32, addr uint64) Instr {
	return MulwDW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulwDW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulw.d.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulwDW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mulw.d.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulwDW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mulw.d.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
