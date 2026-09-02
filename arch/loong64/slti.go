package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slti - slti rd, rj, si12 (2RI12): rd = (rj <s si12) ? 1 : 0.
type Slti struct {
	base

	rd, rj uint8
	imm    imm
}

// Slti - slti rd, rj, si12.
func (Builder) Slti(rd, rj Reg, v Imm12) Instr {
	return Slti{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeSlti(w uint32, addr uint64) Instr {
	return Slti{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i Slti) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slti %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Slti) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["slti"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Slti) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"slti",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
