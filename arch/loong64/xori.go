package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xori - xori rd, rj, ui12 (2RI12): rd = rj ^ ui12.
type Xori struct {
	base

	rd, rj uint8
	imm    imm
}

// Xori - xori rd, rj, ui12.
func (Builder) Xori(rd, rj Reg, v UImm12) Instr {
	return Xori{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeXori(w uint32, addr uint64) Instr {
	return Xori{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 12))),
	}
}

func (i Xori) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("xori %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Xori) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["xori"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Xori) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"xori",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
