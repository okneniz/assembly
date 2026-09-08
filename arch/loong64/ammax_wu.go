package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxWu - ammax.wu rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned.
type AmmaxWu struct {
	base

	rd, rk, rj uint8
}

// AmmaxWu - ammax.wu rd, rk, rj.
func (Builder) AmmaxWu(rd, rk, rj Reg) Instr {
	return AmmaxWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxWu(w uint32) Instr {
	return AmmaxWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax.wu"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
