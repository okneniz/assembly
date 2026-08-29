package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModWu - mod.wu rd, rj, rk (3R): rd = rj % rk (unsigned).
type ModWu struct {
	base

	rd, rj, rk uint8
}

// NewModWu - mod.wu rd, rj, rk.
func NewModWu(rd, rj, rk Reg) Instr {
	return ModWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeModWu(w uint32, addr uint64) Instr {
	return ModWu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i ModWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i ModWu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mod.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i ModWu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mod.wu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
