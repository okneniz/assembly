package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdgtH - ldgt.h rd, rj, rk (DJK): load a half with a bounds check against
// rk, trapping outside.
type LdgtH struct {
	base

	rd, rj, rk uint8
}

// NewLdgtH - ldgt.h rd, rj, rk.
func NewLdgtH(rd, rj, rk Reg) Instr {
	return LdgtH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdgtH(w uint32, addr uint64) Instr {
	return LdgtH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdgtH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldgt.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdgtH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldgt.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdgtH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldgt.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
