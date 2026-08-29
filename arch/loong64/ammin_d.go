package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminD - ammin.d rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): signed fetch-min, old value to rd.
type AmminD struct {
	base

	rd, rk, rj uint8
}

// NewAmminD - ammin.d rd, rk, rj (the operand order the assembler writes).
func NewAmminD(rd, rk, rj Reg) Instr {
	return AmminD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmminD(w uint32, addr uint64) Instr {
	return AmminD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmminD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammin.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmminD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammin.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
