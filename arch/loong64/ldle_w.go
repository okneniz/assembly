package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdleW - ldle.w rd, rj, rk (DJK): load a word with a bounds check against
// rk, trapping outside.
type LdleW struct {
	base

	rd, rj, rk uint8
}

// LdleW - ldle.w rd, rj, rk.
func (Builder) LdleW(rd, rj, rk Reg) Instr {
	return LdleW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdleW(w uint32, addr uint64) Instr {
	return LdleW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdleW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldle.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdleW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldle.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdleW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldle.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
