package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lu52iD - lu52i.d rd, rj, si12 (2RI12): rd = (rj & low52) | (si12 << 52).
type Lu52iD struct {
	base

	rd, rj uint8
	imm    imm
}

// NewLu52iD - lu52i.d rd, rj, si12.
func NewLu52iD(rd, rj Reg, v Imm12) Instr {
	return Lu52iD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeLu52iD(w uint32, addr uint64) Instr {
	return Lu52iD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i Lu52iD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lu52i.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Lu52iD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["lu52i.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Lu52iD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lu52i.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
