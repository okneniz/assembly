package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdgtW - ldgt.w rd, rj, rk (DJK): load a word with a bounds check against
// rk, trapping outside.
type LdgtW struct {
	base

	rd, rj, rk uint8
}

// LdgtW - ldgt.w rd, rj, rk.
func (Builder) LdgtW(rd, rj, rk Reg) Instr {
	return LdgtW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdgtW(w uint32, addr uint64) Instr {
	return LdgtW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdgtW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldgt.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdgtW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldgt.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdgtW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldgt.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
