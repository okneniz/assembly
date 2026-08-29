package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StxD - stx.d rd, rj, rk (3R): MEM[rj + rk] = rd.
type StxD struct {
	base

	rd, rj, rk uint8
}

// NewStxD - stx.d rd, rj, rk.
func NewStxD(rd, rj, rk Reg) Instr {
	return StxD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStxD(w uint32, addr uint64) Instr {
	return StxD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StxD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stx.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StxD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stx.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StxD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stx.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
