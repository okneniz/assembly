package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdgtB - ldgt.b rd, rj, rk (DJK): load a byte with a bounds check against
// rk, trapping outside.
type LdgtB struct {
	base

	rd, rj, rk uint8
}

// NewLdgtB - ldgt.b rd, rj, rk.
func NewLdgtB(rd, rj, rk Reg) Instr {
	return LdgtB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdgtB(w uint32, addr uint64) Instr {
	return LdgtB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdgtB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldgt.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdgtB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldgt.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdgtB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldgt.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
