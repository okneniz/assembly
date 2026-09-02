package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SraW - sra.w rd, rj, rk (3R): rd = sign32(rj >>a (rk & 31)).
type SraW struct {
	base

	rd, rj, rk uint8
}

// SraW - sra.w rd, rj, rk.
func (Builder) SraW(rd, rj, rk Reg) Instr {
	return SraW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSraW(w uint32, addr uint64) Instr {
	return SraW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SraW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sra.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SraW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sra.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SraW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sra.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
