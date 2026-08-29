package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdleD - ldle.d rd, rj, rk (DJK): load a double word with a bounds check
// against rk, trapping outside.
type LdleD struct {
	base

	rd, rj, rk uint8
}

// NewLdleD - ldle.d rd, rj, rk.
func NewLdleD(rd, rj, rk Reg) Instr {
	return LdleD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdleD(w uint32, addr uint64) Instr {
	return LdleD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdleD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldle.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdleD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldle.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdleD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldle.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
