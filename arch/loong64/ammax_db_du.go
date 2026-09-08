package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbDu - ammax_db.du rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = max(MEM[rj], rk), unsigned. A built-in barrier.
type AmmaxDbDu struct {
	base

	rd, rk, rj uint8
}

// AmmaxDbDu - ammax_db.du rd, rk, rj.
func (Builder) AmmaxDbDu(rd, rk, rj Reg) Instr {
	return AmmaxDbDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxDbDu(w uint32) Instr {
	return AmmaxDbDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxDbDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.du %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbDu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammax_db.du"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
