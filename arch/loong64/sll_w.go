package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SllW - sll.w rd, rj, rk (3R): rd = sign32(rj << (rk & 31)).
type SllW struct {
	base

	rd, rj, rk uint8
}

// SllW - sll.w rd, rj, rk.
func (Builder) SllW(rd, rj, rk Reg) Instr {
	return SllW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSllW(w uint32) Instr {
	return SllW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SllW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sll.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SllW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sll.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
