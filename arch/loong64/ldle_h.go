package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdleH - ldle.h rd, rj, rk (DJK): load a half with a bounds check against
// rk, trapping outside.
type LdleH struct {
	base

	rd, rj, rk uint8
}

// NewLdleH - ldle.h rd, rj, rk.
func NewLdleH(rd, rj, rk Reg) Instr {
	return LdleH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdleH(w uint32, addr uint64) Instr {
	return LdleH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdleH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldle.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdleH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldle.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdleH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldle.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
