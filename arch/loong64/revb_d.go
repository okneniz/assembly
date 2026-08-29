package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevbD - revb.d rd, rj (2R): rd = rj with all eight bytes reversed.
type RevbD struct {
	base

	rd, rj uint8
}

// NewRevbD - revb.d rd, rj.
func NewRevbD(rd, rj Reg) Instr {
	return RevbD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevbD(w uint32, addr uint64) Instr {
	return RevbD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RevbD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevbD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["revb.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RevbD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"revb.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
