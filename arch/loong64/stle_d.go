package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StleD - stle.d rd, rj, rk (DJK): store the double word of rd with a
// bounds check against rk, trapping outside.
type StleD struct {
	base

	rd, rj, rk uint8
}

// StleD - stle.d rd, rj, rk.
func (Builder) StleD(rd, rj, rk Reg) Instr {
	return StleD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStleD(w uint32, addr uint64) Instr {
	return StleD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StleD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stle.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StleD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stle.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StleD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stle.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
