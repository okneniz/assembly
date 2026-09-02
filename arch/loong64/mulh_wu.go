package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulhWu - mulh.wu rd, rj, rk (3R): rd = high half of the unsigned 64-bit product low32(rj) * low32(rk).
type MulhWu struct {
	base

	rd, rj, rk uint8
}

// MulhWu - mulh.wu rd, rj, rk.
func (Builder) MulhWu(rd, rj, rk Reg) Instr {
	return MulhWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMulhWu(w uint32, addr uint64) Instr {
	return MulhWu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i MulhWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i MulhWu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mulh.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i MulhWu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mulh.wu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
