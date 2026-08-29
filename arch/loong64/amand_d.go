package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmandD - amand.d rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): fetch-and, old value to rd.
type AmandD struct {
	base

	rd, rk, rj uint8
}

// NewAmandD - amand.d rd, rk, rj (the operand order the assembler writes).
func NewAmandD(rd, rk, rj Reg) Instr {
	return AmandD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmandD(w uint32, addr uint64) Instr {
	return AmandD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmandD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amand.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmandD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amand.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmandD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amand.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
