package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmcasDbH - amcas_db.h rd, rk, rj (3R): if MEM[rj] == rd then MEM[rj] = rk; rd = old MEM[rj]. A built-in barrier.
type AmcasDbH struct {
	base

	rd, rk, rj uint8
}

// AmcasDbH - amcas_db.h rd, rk, rj.
func (Builder) AmcasDbH(rd, rk, rj Reg) Instr {
	return AmcasDbH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmcasDbH(w uint32) Instr {
	return AmcasDbH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmcasDbH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amcas_db.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmcasDbH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amcas_db.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
