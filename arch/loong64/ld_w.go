package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdW - ld.w rd, rj, si12 (2RI12): rd = sign32(MEM[rj + si12]) (the
// offset is an unscaled byte offset).
type LdW struct {
	base

	rd, rj uint8
	imm    imm
}

// LdW - ld.w rd, rj, si12.
func (Builder) LdW(rd, rj Reg, v Imm12) Instr {
	return LdW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeLdW(w uint32, addr uint64) Instr {
	return LdW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i LdW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i LdW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ld.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i LdW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ld.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
