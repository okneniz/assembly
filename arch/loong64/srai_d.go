package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SraiD - srai.d rd, rj, ui6 (2RI6): rd = rj >> (ui6 & 63) arithmetic.
type SraiD struct {
	base

	rd, rj uint8
	imm    imm
}

// NewSraiD - srai.d rd, rj, ui6.
func NewSraiD(rd, rj Reg, v UImm6) Instr {
	return SraiD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeSraiD(w uint32, addr uint64) Instr {
	return SraiD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func (i SraiD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srai.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i SraiD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["srai.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 6)

	return writeWord(w, word)
}

func (i SraiD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srai.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
