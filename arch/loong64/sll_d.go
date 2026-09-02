package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SllD - sll.d rd, rj, rk (3R): rd = rj << (rk & 63).
type SllD struct {
	base

	rd, rj, rk uint8
}

// SllD - sll.d rd, rj, rk.
func (Builder) SllD(rd, rj, rk Reg) Instr {
	return SllD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSllD(w uint32, addr uint64) Instr {
	return SllD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SllD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sll.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SllD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sll.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SllD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sll.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
