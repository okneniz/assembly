package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revbit4B - bitrev.4b rd, rj (2R): rd = rj with the bits reversed in each nibble.
type Revbit4B struct {
	base

	rd, rj uint8
}

// Revbit4B - bitrev.4b rd, rj.
func (Builder) Revbit4B(rd, rj Reg) Instr {
	return Revbit4B{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevbit4B(w uint32, addr uint64) Instr {
	return Revbit4B{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Revbit4B) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.4b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revbit4B) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bitrev.4b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Revbit4B) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bitrev.4b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
