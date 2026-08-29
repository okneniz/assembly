package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtW - stgt.w rd, rj, rk (DJK): store the low word of rd with a bounds
// check against rk, trapping outside.
type StgtW struct {
	base

	rd, rj, rk uint8
}

// NewStgtW - stgt.w rd, rj, rk.
func NewStgtW(rd, rj, rk Reg) Instr {
	return StgtW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStgtW(w uint32, addr uint64) Instr {
	return StgtW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StgtW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stgt.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StgtW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stgt.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
