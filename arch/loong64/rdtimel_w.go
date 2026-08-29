package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RdtimelW - rdtimel.w rd, rj (2R): rd = low 32 bits of the stable counter + rj.
type RdtimelW struct {
	base

	rd, rj uint8
}

// NewRdtimelW - rdtimel.w rd, rj.
func NewRdtimelW(rd, rj Reg) Instr {
	return RdtimelW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRdtimelW(w uint32, addr uint64) Instr {
	return RdtimelW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RdtimelW) ObjDump(_ disasm.ViewCtx) string {
	if i.rd == 0 {
		return "rdcntid.w " + laRegName(i.rj)
	}

	if i.rj == 0 {
		return "rdcntvl.w " + laRegName(i.rd)
	}

	return fmt.Sprintf("rdtimel.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RdtimelW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["rdtimel.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RdtimelW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"rdtimel.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
