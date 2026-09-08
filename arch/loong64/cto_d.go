package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtoD - cto.d rd, rj (2R): rd = the count of trailing ones of rj.
type CtoD struct {
	base

	rd, rj uint8
}

// CtoD - cto.d rd, rj.
func (Builder) CtoD(rd, rj Reg) Instr {
	return CtoD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCtoD(w uint32) Instr {
	return CtoD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CtoD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cto.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtoD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["cto.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
