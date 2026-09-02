package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubW - sub.w rd, rj, rk (3R): rd = sign32(rj - rk).
type SubW struct {
	base

	rd, rj, rk uint8
}

// SubW - sub.w rd, rj, rk.
func (Builder) SubW(rd, rj, rk Reg) Instr {
	return SubW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSubW(w uint32, addr uint64) Instr {
	return SubW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SubW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sub.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SubW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sub.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SubW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sub.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
