package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Maskeqz - maskeqz rd, rj, rk (3R): rd = (rk == 0) ? rj : 0.
type Maskeqz struct {
	base

	rd, rj, rk uint8
}

// NewMaskeqz - maskeqz rd, rj, rk.
func NewMaskeqz(rd, rj, rk Reg) Instr {
	return Maskeqz{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMaskeqz(w uint32, addr uint64) Instr {
	return Maskeqz{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Maskeqz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("maskeqz %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Maskeqz) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["maskeqz"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Maskeqz) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"maskeqz",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
