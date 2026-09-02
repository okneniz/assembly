package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModW - mod.w rd, rj, rk (3R): rd = rj % rk (signed).
type ModW struct {
	base

	rd, rj, rk uint8
}

// ModW - mod.w rd, rj, rk.
func (Builder) ModW(rd, rj, rk Reg) Instr {
	return ModW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeModW(w uint32, addr uint64) Instr {
	return ModW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i ModW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i ModW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mod.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i ModW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mod.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
