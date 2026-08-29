package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ClzD - clz.d rd, rj (2R): rd = the count of leading zeros of rj.
type ClzD struct {
	base

	rd, rj uint8
}

// NewClzD - clz.d rd, rj.
func NewClzD(rd, rj Reg) Instr {
	return ClzD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeClzD(w uint32, addr uint64) Instr {
	return ClzD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i ClzD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clz.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ClzD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["clz.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i ClzD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"clz.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
