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

// NewCtzW - ctz.w rd, rj.
func NewCtzW(rd, rj Reg) Instr {
	return CtzW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCtzW(w uint32, addr uint64) Instr {
	return CtzW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CtzW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ctz.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtzW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ctz.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i CtzW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ctz.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
