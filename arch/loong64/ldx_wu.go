package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdxWu - ldx.wu rd, rj, rk (3R): rd = zero32(MEM[rj + rk]).
type LdxWu struct {
	base

	rd, rj, rk uint8
}

// LdxWu - ldx.wu rd, rj, rk.
func (Builder) LdxWu(rd, rj, rk Reg) Instr {
	return LdxWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeLdxWu(w uint32, addr uint64) Instr {
	return LdxWu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i LdxWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldx.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i LdxWu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ldx.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i LdxWu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ldx.wu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
