package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LlacqD - llacq.d rd, rj (DJ): load-acquire rd = MEM[rj] (64 bits).
type LlacqD struct {
	base

	rd, rj uint8
}

// NewLlacqD - llacq.d rd, rj.
func NewLlacqD(rd, rj Reg) Instr {
	return LlacqD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeLlacqD(w uint32, addr uint64) Instr {
	return LlacqD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i LlacqD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("llacq.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i LlacqD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["llacq.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i LlacqD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"llacq.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
