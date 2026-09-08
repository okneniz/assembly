package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddiD - addi.d rd, rj, si12 (2RI12): rd = rj + si12.
type AddiD struct {
	base

	rd, rj uint8
	imm    imm
}

// AddiD - addi.d rd, rj, si12.
func (Builder) AddiD(rd, rj Reg, v Imm12) Instr {
	return AddiD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeAddiD(w uint32) Instr {
	return AddiD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i AddiD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addi.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i AddiD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["addi.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}
