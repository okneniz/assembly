package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slt - slt rd, rj, rk (3R): rd = (signed rj < signed rk) ? 1 : 0.
type Slt struct {
	base

	rd, rj, rk uint8
}

// Slt - slt rd, rj, rk.
func (Builder) Slt(rd, rj, rk Reg) Instr {
	return Slt{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSlt(w uint32, addr uint64) Instr {
	return Slt{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Slt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slt %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Slt) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["slt"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Slt) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"slt",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
