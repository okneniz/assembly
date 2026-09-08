package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminD - ammin.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), signed.
type AmminD struct {
	base

	rd, rk, rj uint8
}

// AmminD - ammin.d rd, rk, rj.
func (Builder) AmminD(rd, rk, rj Reg) Instr {
	return AmminD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmminD(w uint32) Instr {
	return AmminD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmminD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ammin.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
