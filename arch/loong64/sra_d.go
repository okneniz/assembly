package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SraD - sra.d rd, rj, rk (3R): rd = rj >>a (rk & 63).
type SraD struct {
	base

	rd, rj, rk uint8
}

// NewSraD - sra.d rd, rj, rk.
func NewSraD(rd, rj, rk Reg) Instr {
	return SraD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSraD(w uint32, addr uint64) Instr {
	return SraD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SraD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sra.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SraD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sra.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SraD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sra.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
