package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SllW - sll.w rd, rj, rk (3R): rd = sign32(rj << (rk & 31)).
type SllW struct {
	base

	rd, rj, rk uint8
}

// NewSllW - sll.w rd, rj, rk.
func NewSllW(rd, rj, rk Reg) Instr {
	return SllW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSllW(w uint32, addr uint64) Instr {
	return SllW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SllW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sll.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SllW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sll.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SllW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sll.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
