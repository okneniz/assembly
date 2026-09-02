package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revh2W - revh.2w rd, rj (2R): rd = rj with the halfwords reversed in each word.
type Revh2W struct {
	base

	rd, rj uint8
}

// Revh2W - revh.2w rd, rj.
func (Builder) Revh2W(rd, rj Reg) Instr {
	return Revh2W{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevh2W(w uint32, addr uint64) Instr {
	return Revh2W{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Revh2W) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revh.2w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revh2W) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["revh.2w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Revh2W) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"revh.2w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
