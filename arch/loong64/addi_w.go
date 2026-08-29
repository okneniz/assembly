package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddiW - addi.w rd, rj, si12 (2RI12): rd = sign32(rj + si12).
type AddiW struct {
	base

	rd, rj uint8
	imm    imm
}

// NewAddiW - addi.w rd, rj, si12.
func NewAddiW(rd, rj Reg, v Imm12) Instr {
	return AddiW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeAddiW(w uint32, addr uint64) Instr {
	return AddiW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i AddiW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addi.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i AddiW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["addi.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i AddiW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"addi.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
