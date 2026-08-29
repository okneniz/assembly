package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWDW - crcc.w.d.w rd, rj, rk (DJK): carryless CRC32 of rj (double) into
// rk's accumulator.
type CrccWDW struct {
	base

	rd, rj, rk uint8
}

// NewCrccWDW - crcc.w.d.w rd, rj, rk.
func NewCrccWDW(rd, rj, rk Reg) Instr {
	return CrccWDW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrccWDW(w uint32, addr uint64) Instr {
	return CrccWDW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrccWDW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.d.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWDW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crcc.w.d.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrccWDW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crcc.w.d.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
