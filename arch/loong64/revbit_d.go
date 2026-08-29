package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevbitD - bitrev.d rd, rj (2R): rd = rj with all 64 bits reversed.
type RevbitD struct {
	base

	rd, rj uint8
}

// NewRevbitD - bitrev.d rd, rj.
func NewRevbitD(rd, rj Reg) Instr {
	return RevbitD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevbitD(w uint32, addr uint64) Instr {
	return RevbitD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RevbitD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevbitD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bitrev.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RevbitD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bitrev.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
