package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtD - stgt.d rd, rj, rk (DJK): store the double word of rd with a
// bounds check against rk, trapping outside.
type StgtD struct {
	base

	rd, rj, rk uint8
}

// NewStgtD - stgt.d rd, rj, rk.
func NewStgtD(rd, rj, rk Reg) Instr {
	return StgtD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStgtD(w uint32, addr uint64) Instr {
	return StgtD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StgtD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["stgt.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i StgtD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"stgt.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
