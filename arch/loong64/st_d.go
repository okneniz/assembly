package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StD - st.d rd, rj, si12 (2RI12): MEM[rj + si12] = rd (the offset is
// an unscaled byte offset).
type StD struct {
	base

	rd, rj uint8
	imm    imm
}

// StD - st.d rd, rj, si12.
func (Builder) StD(rd, rj Reg, v Imm12) Instr {
	return StD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeStD(w uint32, addr uint64) Instr {
	return StD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func (i StD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("st.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i StD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["st.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 12)

	return writeWord(w, word)
}

func (i StD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"st.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
