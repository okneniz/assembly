package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWBW - crcc.w.b.w rd, rj, rk (DJK): carryless CRC32 of rj (byte) into
// rk's accumulator.
type CrccWBW struct {
	base

	rd, rj, rk uint8
}

// CrccWBW - crcc.w.b.w rd, rj, rk.
func (Builder) CrccWBW(rd, rj, rk Reg) Instr {
	return CrccWBW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrccWBW(w uint32, addr uint64) Instr {
	return CrccWBW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrccWBW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.b.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWBW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crcc.w.b.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrccWBW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crcc.w.b.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
