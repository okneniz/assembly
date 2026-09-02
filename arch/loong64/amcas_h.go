package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasH - amcas.h rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj].
type AmcasH struct {
	base

	rd, rk, rj uint8
}

// AmcasH - amcas.h rd, rk, rj.
func (Builder) AmcasH(rd, rk, rj Reg) Instr {
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
