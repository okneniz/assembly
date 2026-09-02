package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StleW - stle.w rd, rj, rk (DJK): store the low word of rd with a bounds
// check against rk, trapping outside.
type StleW struct {
	base

	rd, rj, rk uint8
}

// StleW - stle.w rd, rj, rk.
func (Builder) StleW(rd, rj, rk Reg) Instr {
	return StleW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStleW(w uint32, addr uint64) Instr {
	return StleW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StleW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stle.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StleW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stle.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StleW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stle.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
