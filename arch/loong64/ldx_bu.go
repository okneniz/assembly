package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxBu - ldx.bu rd, rj, rk (3R): rd = the byte at rj + rk, zero-extended.
type LdxBu struct {
	base

	rd, rj, rk uint8
}

// NewLdxBu - ldx.bu rd, rj, rk.
func NewLdxBu(rd, rj, rk Reg) Instr {
	return LdxBu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdxBu(w uint32, addr uint64) Instr {
	return LdxBu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdxBu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.bu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxBu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldx.bu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxBu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldx.bu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
