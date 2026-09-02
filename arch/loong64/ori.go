package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ori - ori rd, rj, ui12 (2RI12): rd = rj | ui12.
type Ori struct {
	base

	rd, rj uint8
	imm    imm
}

// Ori - ori rd, rj, ui12.
func (Builder) Ori(rd, rj Reg, v UImm12) Instr {
	return Ori{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeOri(w uint32, addr uint64) Instr {
	return Ori{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 12))),
	}
}

func (i Ori) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ori %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Ori) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ori"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i Ori) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ori",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
