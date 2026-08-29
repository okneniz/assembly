package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxHu - ldx.hu rd, rj, rk (3R): rd = the halfword at rj + rk, zero-extended.
type LdxHu struct {
	base

	rd, rj, rk uint8
}

// NewLdxHu - ldx.hu rd, rj, rk.
func NewLdxHu(rd, rj, rk Reg) Instr {
	return LdxHu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdxHu(w uint32, addr uint64) Instr {
	return LdxHu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdxHu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.hu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxHu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldx.hu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxHu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldx.hu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
