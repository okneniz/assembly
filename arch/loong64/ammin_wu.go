package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminWu - ammin.wu rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] = min(MEM[rj], rk), unsigned.
type AmminWu struct {
	base

	rd, rk, rj uint8
}

// AmminWu - ammin.wu rd, rk, rj.
func (Builder) AmminWu(rd, rk, rj Reg) Instr {
	return AmminWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmminWu(w uint32, addr uint64) Instr {
	return AmminWu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmminWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminWu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammin.wu"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmminWu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammin.wu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
