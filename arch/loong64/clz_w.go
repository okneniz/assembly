package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ClzW - clz.w rd, rj (2R): rd = the count of leading zeros of low32(rj).
type ClzW struct {
	base

	rd, rj uint8
}

// NewClzW - clz.w rd, rj.
func NewClzW(rd, rj Reg) Instr {
	return ClzW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeClzW(w uint32, addr uint64) Instr {
	return ClzW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i ClzW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clz.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ClzW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["clz.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i ClzW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"clz.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
