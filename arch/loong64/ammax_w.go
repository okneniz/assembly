package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxW - ammax.w rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed.
type AmmaxW struct {
	base

	rd, rk, rj uint8
}

// AmmaxW - ammax.w rd, rk, rj.
func (Builder) AmmaxW(rd, rk, rj Reg) Instr {
	return AmmaxW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxW(w uint32, addr uint64) Instr {
	return AmmaxW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.w %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammax.w"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmmaxW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammax.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
