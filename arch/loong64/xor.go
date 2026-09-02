package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xor - xor rd, rj, rk (3R): rd = rj ^ rk.
type Xor struct {
	base

	rd, rj, rk uint8
}

// Xor - xor rd, rj, rk.
func (Builder) Xor(rd, rj, rk Reg) Instr {
	return Xor{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeXor(w uint32, addr uint64) Instr {
	return Xor{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i Xor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("xor %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i Xor) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["xor"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i Xor) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"xor",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
