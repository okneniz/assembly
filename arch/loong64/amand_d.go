package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmandD - amand.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] &= rk.
type AmandD struct {
	base

	rd, rk, rj uint8
}

// AmandD - amand.d rd, rk, rj.
func (Builder) AmandD(rd, rk, rj Reg) Instr {
	return AmandD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmandD(w uint32) Instr {
	return AmandD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmandD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amand.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmandD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amand.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
