package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DivW - div.w rd, rj, rk (3R): rd = rj / rk (signed).
type DivW struct {
	base

	rd, rj, rk uint8
}

// DivW - div.w rd, rj, rk.
func (Builder) DivW(rd, rj, rk Reg) Instr {
	return DivW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeDivW(w uint32, addr uint64) Instr {
	return DivW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i DivW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i DivW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["div.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i DivW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"div.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
