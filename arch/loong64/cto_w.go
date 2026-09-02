package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtoW - cto.w rd, rj (2R): rd = the count of trailing ones of low32(rj).
type CtoW struct {
	base

	rd, rj uint8
}

// CtoW - cto.w rd, rj.
func (Builder) CtoW(rd, rj Reg) Instr {
	return CtoW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCtoW(w uint32, addr uint64) Instr {
	return CtoW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CtoW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cto.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtoW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["cto.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i CtoW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"cto.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
