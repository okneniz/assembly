package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DivD - div.d rd, rj, rk (3R): rd = rj / rk (signed).
type DivD struct {
	base

	rd, rj, rk uint8
}

// NewDivD - div.d rd, rj, rk.
func NewDivD(rd, rj, rk Reg) Instr {
	return DivD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeDivD(w uint32, addr uint64) Instr {
	return DivD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i DivD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i DivD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["div.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i DivD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"div.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
