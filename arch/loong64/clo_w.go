package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CloW - clo.w rd, rj (2R): rd = the count of leading ones of low32(rj).
type CloW struct {
	base

	rd, rj uint8
}

// CloW - clo.w rd, rj.
func (Builder) CloW(rd, rj Reg) Instr {
	return CloW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCloW(w uint32) Instr {
	return CloW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CloW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clo.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CloW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["clo.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
