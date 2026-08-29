package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LlacqW - llacq.w rd, rj (DJ): load-acquire rd = MEM[rj] (32 bits).
type LlacqW struct {
	base

	rd, rj uint8
}

// NewLlacqW - llacq.w rd, rj.
func NewLlacqW(rd, rj Reg) Instr {
	return LlacqW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeLlacqW(w uint32, addr uint64) Instr {
	return LlacqW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i LlacqW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("llacq.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i LlacqW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["llacq.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i LlacqW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"llacq.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
