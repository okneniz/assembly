package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cpucfg - cpucfg rd, rj (2R): rd = the configuration register selected by rj.
type Cpucfg struct {
	base

	rd, rj uint8
}

// Cpucfg - cpucfg rd, rj.
func (Builder) Cpucfg(rd, rj Reg) Instr {
	return Cpucfg{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCpucfg(w uint32, addr uint64) Instr {
	return Cpucfg{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Cpucfg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cpucfg %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Cpucfg) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["cpucfg"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Cpucfg) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"cpucfg",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
