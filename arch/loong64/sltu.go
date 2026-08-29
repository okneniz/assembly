package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sltu - sltu rd, rj, rk (3R): rd = (rj < rk) ? 1 : 0 (unsigned).
type Sltu struct {
	base

	rd, rj, rk uint8
}

// NewSltu - sltu rd, rj, rk.
func NewSltu(rd, rj, rk Reg) Instr {
	return Sltu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSltu(w uint32, addr uint64) Instr {
	return Sltu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Sltu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sltu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Sltu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sltu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Sltu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sltu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
