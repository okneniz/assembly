package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DivWu - div.wu rd, rj, rk (3R): rd = rj / rk (unsigned).
type DivWu struct {
	base

	rd, rj, rk uint8
}

// DivWu - div.wu rd, rj, rk.
func (Builder) DivWu(rd, rj, rk Reg) Instr {
	return DivWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeDivWu(w uint32) Instr {
	return DivWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i DivWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i DivWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["div.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
