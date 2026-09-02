package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmaddH - amadd.h rd, rk, rj (3R): rd = old MEM[rj]; MEM[rj] += rk.
type AmaddH struct {
	base

	rd, rk, rj uint8
}

// AmaddH - amadd.h rd, rk, rj.
func (Builder) AmaddH(rd, rk, rj Reg) Instr {
	return AmaddH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmaddH(w uint32, addr uint64) Instr {
	return AmaddH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmaddH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amadd.h %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmaddH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["amadd.h"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmaddH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amadd.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
