package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbWu - ammax_db.wu rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned. A built-in barrier.
type AmmaxDbWu struct {
	base

	rd, rk, rj uint8
}

// AmmaxDbWu - ammax_db.wu rd, rk, rj.
func (Builder) AmmaxDbWu(rd, rk, rj Reg) Instr {
	return AmmaxDbWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxDbWu(w uint32) Instr {
	return AmmaxDbWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxDbWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax_db.wu"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
