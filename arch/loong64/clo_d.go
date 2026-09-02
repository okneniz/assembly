package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CloD - clo.d rd, rj (2R): rd = the count of leading ones of rj.
type CloD struct {
	base

	rd, rj uint8
}

// CloD - clo.d rd, rj.
func (Builder) CloD(rd, rj Reg) Instr {
	return CloD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCloD(w uint32, addr uint64) Instr {
	return CloD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CloD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clo.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CloD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["clo.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i CloD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"clo.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
