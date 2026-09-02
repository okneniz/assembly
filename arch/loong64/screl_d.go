package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScrelD - screl.d rd, rj (DJ): store-release MEM[rj] = rd (64 bits).
type ScrelD struct {
	base

	rd, rj uint8
}

// ScrelD - screl.d rd, rj.
func (Builder) ScrelD(rd, rj Reg) Instr {
	return ScrelD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeScrelD(w uint32, addr uint64) Instr {
	return ScrelD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i ScrelD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("screl.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i ScrelD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["screl.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i ScrelD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"screl.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
