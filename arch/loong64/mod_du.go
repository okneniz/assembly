package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ModDu - mod.du rd, rj, rk (3R): rd = rj % rk (unsigned).
type ModDu struct {
	base

	rd, rj, rk uint8
}

// ModDu - mod.du rd, rj, rk.
func (Builder) ModDu(rd, rj, rk Reg) Instr {
	return ModDu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeModDu(w uint32, addr uint64) Instr {
	return ModDu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i ModDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mod.du %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i ModDu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["mod.du"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i ModDu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mod.du",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
