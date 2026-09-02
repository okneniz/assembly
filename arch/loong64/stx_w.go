package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StxW - stx.w rd, rj, rk (3R): MEM[rj + rk] = low32(rd).
type StxW struct {
	base

	rd, rj, rk uint8
}

// StxW - stx.w rd, rj, rk.
func (Builder) StxW(rd, rj, rk Reg) Instr {
	return StxW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStxW(w uint32, addr uint64) Instr {
	return StxW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StxW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stx.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StxW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stx.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StxW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stx.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
