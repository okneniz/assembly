package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxW - ldx.w rd, rj, rk (3R): rd = sign32(MEM[rj + rk]).
type LdxW struct {
	base

	rd, rj, rk uint8
}

// LdxW - ldx.w rd, rj, rk.
func (Builder) LdxW(rd, rj, rk Reg) Instr {
	return LdxW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdxW(w uint32, addr uint64) Instr {
	return LdxW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdxW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldx.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldx.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
