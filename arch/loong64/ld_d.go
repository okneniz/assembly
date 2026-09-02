package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdD - ld.d rd, rj, si12 (2RI12): rd = MEM[rj + si12] (the offset is
// an unscaled byte offset).
type LdD struct {
	base

	rd, rj uint8
	imm    imm
}

// LdD - ld.d rd, rj, si12.
func (Builder) LdD(rd, rj Reg, v Imm12) Instr {
	return LdD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeLdD(w uint32, addr uint64) Instr {
	return LdD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i LdD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i LdD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ld.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i LdD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ld.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
