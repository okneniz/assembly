package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmaddD - amadd.d rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] += rk.
type AmaddD struct {
	base

	rd, rk, rj uint8
}

// AmaddD - amadd.d rd, rk, rj.
func (Builder) AmaddD(rd, rk, rj Reg) Instr {
	return AmaddD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmaddD(w uint32) Instr {
	return AmaddD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmaddD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amadd.d %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmaddD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["amadd.d"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}
