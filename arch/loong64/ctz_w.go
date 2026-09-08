package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtzW - ctz.w rd, rj (2R): rd = the count of trailing zeros of low32(rj).
type CtzW struct {
	base

	rd, rj uint8
}

// CtzW - ctz.w rd, rj.
func (Builder) CtzW(rd, rj Reg) Instr {
	return CtzW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCtzW(w uint32) Instr {
	return CtzW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CtzW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ctz.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtzW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ctz.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
