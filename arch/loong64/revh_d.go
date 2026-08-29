package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevhD - revh.d rd, rj (2R): rd = rj with all four halfwords reversed.
type RevhD struct {
	base

	rd, rj uint8
}

// NewRevhD - revh.d rd, rj.
func NewRevhD(rd, rj Reg) Instr {
	return RevhD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevhD(w uint32, addr uint64) Instr {
	return RevhD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RevhD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revh.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevhD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["revh.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RevhD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"revh.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
