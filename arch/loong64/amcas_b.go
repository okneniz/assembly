package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasB - amcas.b rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj].
type AmcasB struct {
	base

	rd, rk, rj uint8
}

// AmcasB - amcas.b rd, rk, rj.
func (Builder) AmcasB(rd, rk, rj Reg) Instr {
	return AmcasB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmcasB(w uint32, addr uint64) Instr {
	return AmcasB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmcasB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas.b %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amcas.b"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmcasB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amcas.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
