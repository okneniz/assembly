package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdleB - ldle.b rd, rj, rk (DJK): load a byte with a bounds check against
// rk, trapping outside.
type LdleB struct {
	base

	rd, rj, rk uint8
}

// LdleB - ldle.b rd, rj, rk.
func (Builder) LdleB(rd, rj, rk Reg) Instr {
	return LdleB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdleB(w uint32, addr uint64) Instr {
	return LdleB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdleB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldle.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdleB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldle.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdleB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldle.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
