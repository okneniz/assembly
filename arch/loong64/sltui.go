package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sltui - sltui rd, rj, si12 (2RI12): rd = (rj <u sext32(si12)) ? 1 : 0.
type Sltui struct {
	base

	rd, rj uint8
	imm    imm
}

// NewSltui - sltui rd, rj, si12.
func NewSltui(rd, rj Reg, v Imm12) Instr {
	return Sltui{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeSltui(w uint32, addr uint64) Instr {
	return Sltui{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i Sltui) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sltui %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Sltui) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["sltui"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Sltui) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sltui",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
