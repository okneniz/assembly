package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasH - amcas.h rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): compare-swap, old value to rd.
type AmcasH struct {
	base

	rd, rk, rj uint8
}

// NewAmcasH - amcas.h rd, rk, rj (the operand order the assembler writes).
func NewAmcasH(rd, rk, rj Reg) Instr {
	return AmcasH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmcasH(w uint32, addr uint64) Instr {
	return AmcasH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmcasH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amcas.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmcasH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amcas.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
