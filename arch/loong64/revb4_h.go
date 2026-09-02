package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revb4H - revb.4h rd, rj (2R): rd = rj with the bytes reversed in each halfword.
type Revb4H struct {
	base

	rd, rj uint8
}

// Revb4H - revb.4h rd, rj.
func (Builder) Revb4H(rd, rj Reg) Instr {
	return Revb4H{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevb4H(w uint32, addr uint64) Instr {
	return Revb4H{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Revb4H) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.4h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revb4H) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["revb.4h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Revb4H) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"revb.4h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
