package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SlliD - slli.d rd, rj, ui6 (2RI6): rd = rj << (ui6 & 63).
type SlliD struct {
	base

	rd, rj uint8
	imm    imm
}

// SlliD - slli.d rd, rj, ui6.
func (Builder) SlliD(rd, rj Reg, v UImm6) Instr {
	return SlliD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeSlliD(w uint32, addr uint64) Instr {
	return SlliD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func (i SlliD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slli.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i SlliD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["slli.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 6)

	return writeWord(w, word)
}

func (i SlliD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"slli.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
