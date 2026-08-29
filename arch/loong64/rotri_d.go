package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RotriD - rotri.d rd, rj, ui6 (2RI6): rd = rotate_right64(rj, ui6 & 63).
type RotriD struct {
	base

	rd, rj uint8
	imm    imm
}

// NewRotriD - rotri.d rd, rj, ui6.
func NewRotriD(rd, rj Reg, v UImm6) Instr {
	return RotriD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeRotriD(w uint32, addr uint64) Instr {
	return RotriD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func (i RotriD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rotri.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i RotriD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["rotri.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 6)

	return writeWord(w, word)
}

func (i RotriD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"rotri.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
