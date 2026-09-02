package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxD - ammax.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), signed.
type AmmaxD struct {
	base

	rd, rk, rj uint8
}

// AmmaxD - ammax.d rd, rk, rj.
func (Builder) AmmaxD(rd, rk, rj Reg) Instr {
	return AmmaxD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxD(w uint32, addr uint64) Instr {
	return AmmaxD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammax.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmmaxD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammax.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
