package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrliW - srli.w rd, rj, ui5 (2RI5): rd = sign32(rj >> ui5) logical.
type SrliW struct {
	base

	rd, rj uint8
	imm    imm
}

// SrliW - srli.w rd, rj, ui5.
func (Builder) SrliW(rd, rj Reg, v UImm5) Instr {
	return SrliW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeSrliW(w uint32, addr uint64) Instr {
	return SrliW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 5))),
	}
}

func (i SrliW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srli.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i SrliW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["srli.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.imm.val, 10, 5)

	return writeWord(w, word)
}

func (i SrliW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srli.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
