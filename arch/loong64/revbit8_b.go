package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revbit8B - bitrev.8b rd, rj (2R): rd = rj with the bits reversed in each byte.
type Revbit8B struct {
	base

	rd, rj uint8
}

// Revbit8B - bitrev.8b rd, rj.
func (Builder) Revbit8B(rd, rj Reg) Instr {
	return Revbit8B{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevbit8B(w uint32, addr uint64) Instr {
	return Revbit8B{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Revbit8B) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.8b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revbit8B) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bitrev.8b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Revbit8B) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bitrev.8b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
