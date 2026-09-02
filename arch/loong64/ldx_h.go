package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxH - ldx.h rd, rj, rk (3R): rd = sign16(MEM[rj + rk]).
type LdxH struct {
	base

	rd, rj, rk uint8
}

// LdxH - ldx.h rd, rj, rk.
func (Builder) LdxH(rd, rj, rk Reg) Instr {
	return LdxH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdxH(w uint32, addr uint64) Instr {
	return LdxH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdxH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldx.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldx.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
