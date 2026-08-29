package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddD - add.d rd, rj, rk (3R): rd = rj + rk.
type AddD struct {
	base

	rd, rj, rk uint8
}

// NewAddD - add.d rd, rj, rk.
func NewAddD(rd, rj, rk Reg) Instr {
	return AddD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAddD(w uint32, addr uint64) Instr {
	return AddD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i AddD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i AddD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["add.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i AddD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"add.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
