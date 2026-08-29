package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Andn - andn rd, rj, rk (3R): rd = rj & ~rk.
type Andn struct {
	base

	rd, rj, rk uint8
}

// NewAndn - andn rd, rj, rk.
func NewAndn(rd, rj, rk Reg) Instr {
	return Andn{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAndn(w uint32, addr uint64) Instr {
	return Andn{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Andn) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("andn %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Andn) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["andn"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Andn) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"andn",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
