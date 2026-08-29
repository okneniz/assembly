package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtH - stgt.h rd, rj, rk (DJK): store the low half of rd with a bounds
// check against rk, trapping outside.
type StgtH struct {
	base

	rd, rj, rk uint8
}

// NewStgtH - stgt.h rd, rj, rk.
func NewStgtH(rd, rj, rk Reg) Instr {
	return StgtH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStgtH(w uint32, addr uint64) Instr {
	return StgtH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StgtH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stgt.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StgtH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stgt.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
