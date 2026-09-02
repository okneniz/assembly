package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrccWWW - crcc.w.w.w rd, rj, rk (DJK): carryless CRC32 of rj (word) into
// rk's accumulator.
type CrccWWW struct {
	base

	rd, rj, rk uint8
}

// CrccWWW - crcc.w.w.w rd, rj, rk.
func (Builder) CrccWWW(rd, rj, rk Reg) Instr {
	return CrccWWW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrccWWW(w uint32, addr uint64) Instr {
	return CrccWWW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrccWWW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crcc.w.w.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrccWWW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crcc.w.w.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrccWWW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crcc.w.w.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
