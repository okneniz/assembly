package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Revb2W - revb.2w rd, rj (2R): rd = rj with the bytes reversed in each word.
type Revb2W struct {
	base

	rd, rj uint8
}

// Revb2W - revb.2w rd, rj.
func (Builder) Revb2W(rd, rj Reg) Instr {
	return Revb2W{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevb2W(w uint32, addr uint64) Instr {
	return Revb2W{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i Revb2W) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("revb.2w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i Revb2W) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["revb.2w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i Revb2W) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"revb.2w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
