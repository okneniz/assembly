package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModD - mod.d rd, rj, rk (3R): rd = rj % rk (signed).
type ModD struct {
	base

	rd, rj, rk uint8
}

// NewModD - mod.d rd, rj, rk.
func NewModD(rd, rj, rk Reg) Instr {
	return ModD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeModD(w uint32, addr uint64) Instr {
	return ModD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i ModD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i ModD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mod.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i ModD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mod.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
