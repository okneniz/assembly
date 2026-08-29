package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmorD - amor.d rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): fetch-or, old value to rd.
type AmorD struct {
	base

	rd, rk, rj uint8
}

// NewAmorD - amor.d rd, rk, rj (the operand order the assembler writes).
func NewAmorD(rd, rk, rj Reg) Instr {
	return AmorD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmorD(w uint32, addr uint64) Instr {
	return AmorD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmorD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amor.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmorD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amor.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmorD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amor.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
