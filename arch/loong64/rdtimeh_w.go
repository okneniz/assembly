package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RdtimehW - rdtimeh.w rd, rj (2R): rd = high 32 bits of the stable counter + rj.
type RdtimehW struct {
	base

	rd, rj uint8
}

// NewRdtimehW - rdtimeh.w rd, rj.
func NewRdtimehW(rd, rj Reg) Instr {
	return RdtimehW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRdtimehW(w uint32, addr uint64) Instr {
	return RdtimehW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RdtimehW) ObjDump(_ disasm.ViewCtx) string {
	if i.rj == 0 {
		return "rdcntvh.w " + laRegName(i.rd)
	}

	return fmt.Sprintf("rdtimeh.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RdtimehW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["rdtimeh.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RdtimehW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"rdtimeh.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
