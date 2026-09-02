package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Orn - orn rd, rj, rk (3R): rd = rj | ~rk.
type Orn struct {
	base

	rd, rj, rk uint8
}

// Orn - orn rd, rj, rk.
func (Builder) Orn(rd, rj, rk Reg) Instr {
	return Orn{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeOrn(w uint32, addr uint64) Instr {
	return Orn{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Orn) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("orn %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Orn) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["orn"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Orn) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"orn",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
