package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Nor - nor rd, rj, rk (3R): rd = ~(rj | rk).
type Nor struct {
	base

	rd, rj, rk uint8
}

// Nor - nor rd, rj, rk.
func (Builder) Nor(rd, rj, rk Reg) Instr {
	return Nor{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeNor(w uint32, addr uint64) Instr {
	return Nor{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Nor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("nor %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Nor) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["nor"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Nor) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"nor",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
