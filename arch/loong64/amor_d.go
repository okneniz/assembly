package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmorD - amor.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] |= rk.
type AmorD struct {
	base

	rd, rk, rj uint8
}

// AmorD - amor.d rd, rk, rj.
func (Builder) AmorD(rd, rk, rj Reg) Instr {
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
