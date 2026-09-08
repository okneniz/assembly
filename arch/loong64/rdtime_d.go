package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RdtimeD - rdtime.d rd, rj (2R): rd = the 64-bit stable counter + rj.
type RdtimeD struct {
	base

	rd, rj uint8
}

// RdtimeD - rdtime.d rd, rj.
func (Builder) RdtimeD(rd, rj Reg) Instr {
	return RdtimeD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRdtimeD(w uint32) Instr {
	return RdtimeD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RdtimeD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rdtime.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RdtimeD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["rdtime.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
