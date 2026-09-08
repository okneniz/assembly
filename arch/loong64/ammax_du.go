package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDu - ammax.du rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned.
type AmmaxDu struct {
	base

	rd, rk, rj uint8
}

// AmmaxDu - ammax.du rd, rk, rj.
func (Builder) AmmaxDu(rd, rk, rj Reg) Instr {
	return AmmaxDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxDu(w uint32) Instr {
	return AmmaxDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax.du %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax.du"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
