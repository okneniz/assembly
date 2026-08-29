package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StH - st.h rd, rj, si12 (2RI12): MEM[rj + si12] = low16(rd) (the
// offset is an unscaled byte offset).
type StH struct {
	base

	rd, rj uint8
	imm    imm
}

// NewStH - st.h rd, rj, si12.
func NewStH(rd, rj Reg, v Imm12) Instr {
	return StH{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeStH(w uint32, addr uint64) Instr {
	return StH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i StH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("st.h %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i StH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["st.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i StH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"st.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
