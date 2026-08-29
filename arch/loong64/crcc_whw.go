package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWHW - crcc.w.h.w rd, rj, rk (DJK): carryless CRC32 of rj (half) into
// rk's accumulator.
type CrccWHW struct {
	base

	rd, rj, rk uint8
}

// NewCrccWHW - crcc.w.h.w rd, rj, rk.
func NewCrccWHW(rd, rj, rk Reg) Instr {
	return CrccWHW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrccWHW(w uint32, addr uint64) Instr {
	return CrccWHW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrccWHW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.h.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWHW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crcc.w.h.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrccWHW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crcc.w.h.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
