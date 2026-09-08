package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmaddH - amadd.h rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] += rk.
type AmaddH struct {
	base

	rd, rk, rj uint8
}

// AmaddH - amadd.h rd, rk, rj.
func (Builder) AmaddH(rd, rk, rj Reg) Instr {
	return AmaddH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmaddH(w uint32) Instr {
	return AmaddH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmaddH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amadd.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmaddH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amadd.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
