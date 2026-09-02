package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScrelW - screl.w rd, rj (DJ): store-release MEM[rj] = rd (32 bits).
type ScrelW struct {
	base

	rd, rj uint8
}

// ScrelW - screl.w rd, rj.
func (Builder) ScrelW(rd, rj Reg) Instr {
	return ScrelW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeScrelW(w uint32, addr uint64) Instr {
	return ScrelW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i ScrelW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("screl.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ScrelW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["screl.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i ScrelW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"screl.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
