package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrlW - srl.w rd, rj, rk (3R): rd = sign32(rj >>l (rk & 31)).
type SrlW struct {
	base

	rd, rj, rk uint8
}

// NewSrlW - srl.w rd, rj, rk.
func NewSrlW(rd, rj, rk Reg) Instr {
	return SrlW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSrlW(w uint32, addr uint64) Instr {
	return SrlW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SrlW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srl.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SrlW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["srl.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SrlW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srl.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
