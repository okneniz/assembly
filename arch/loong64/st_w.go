package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StW - st.w rd, rj, si12 (2RI12): MEM[rj + si12] = low32(rd) (the
// offset is an unscaled byte offset).
type StW struct {
	base

	rd, rj uint8
	imm    imm
}

// StW - st.w rd, rj, si12.
func (Builder) StW(rd, rj Reg, v Imm12) Instr {
	return StW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeStW(w uint32, addr uint64) Instr {
	return StW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i StW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("st.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i StW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["st.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i StW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"st.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
