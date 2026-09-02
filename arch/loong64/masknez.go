package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Masknez - masknez rd, rj, rk (3R): rd = (rk != 0) ? rj : 0.
type Masknez struct {
	base

	rd, rj, rk uint8
}

// Masknez - masknez rd, rj, rk.
func (Builder) Masknez(rd, rj, rk Reg) Instr {
	return Masknez{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeMasknez(w uint32, addr uint64) Instr {
	return Masknez{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Masknez) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("masknez %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Masknez) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["masknez"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Masknez) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"masknez",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
