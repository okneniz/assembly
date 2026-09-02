package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// And - and rd, rj, rk (3R): rd = rj & rk.
type And struct {
	base

	rd, rj, rk uint8
}

// And - and rd, rj, rk.
func (Builder) And(rd, rj, rk Reg) Instr {
	return And{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAnd(w uint32, addr uint64) Instr {
	return And{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i And) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("and %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i And) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["and"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i And) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"and",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
