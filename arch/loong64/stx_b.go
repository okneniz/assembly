package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StxB - stx.b rd, rj, rk (3R): MEM[rj + rk] = low8(rd).
type StxB struct {
	base

	rd, rj, rk uint8
}

// StxB - stx.b rd, rj, rk.
func (Builder) StxB(rd, rj, rk Reg) Instr {
	return StxB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStxB(w uint32) Instr {
	return StxB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StxB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stx.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StxB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stx.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
