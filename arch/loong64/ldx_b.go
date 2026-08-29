package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxB - ldx.b rd, rj, rk (3R): rd = sign8(MEM[rj + rk]).
type LdxB struct {
	base

	rd, rj, rk uint8
}

// NewLdxB - ldx.b rd, rj, rk.
func NewLdxB(rd, rj, rk Reg) Instr {
	return LdxB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdxB(w uint32, addr uint64) Instr {
	return LdxB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdxB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldx.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldx.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
