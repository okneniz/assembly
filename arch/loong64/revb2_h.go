package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revb2H - revb.2h rd, rj (2R): rd = rj with the bytes reversed in each halfword of the low 32 bits.
type Revb2H struct {
	base

	rd, rj uint8
}

// Revb2H - revb.2h rd, rj.
func (Builder) Revb2H(rd, rj Reg) Instr {
	return Revb2H{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevb2H(w uint32, addr uint64) Instr {
	return Revb2H{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Revb2H) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.2h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revb2H) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["revb.2h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Revb2H) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"revb.2h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
