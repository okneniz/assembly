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

// NewCloW - clo.w rd, rj.
func NewCloW(rd, rj Reg) Instr {
	return CloW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCloW(w uint32, addr uint64) Instr {
	return CloW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CloW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clo.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CloW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["clo.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i CloW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"clo.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
