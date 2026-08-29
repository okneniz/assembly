package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Or - or rd, rj, rk (3R): rd = rj | rk.
type Or struct {
	base

	rd, rj, rk uint8
}

// NewOr - or rd, rj, rk.
func NewOr(rd, rj, rk Reg) Instr {
	return Or{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeOr(w uint32, addr uint64) Instr {
	return Or{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Or) ObjDump(_ disasm.ViewCtx) string {
	if i.rk == 0 {
		return fmt.Sprintf("move %s, %s", laRegName(i.rd), laRegName(i.rj))
	}

	return fmt.Sprintf("or %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Or) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["or"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Or) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"or",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
