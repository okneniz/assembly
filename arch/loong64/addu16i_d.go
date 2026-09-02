package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Addu16iD - addu16i.d rd, rj, si16 (2RI16): rd = rj + sext(si16 << 16).
// The assembly writes the plain si16 (the << 16 is part of the semantics).
type Addu16iD struct {
	base

	rd, rj uint8
	imm    imm
}

// Addu16iD - addu16i.d rd, rj, si16.
func (Builder) Addu16iD(rd, rj Reg, v Imm16) Instr {
	return Addu16iD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeAddu16iD(w uint32, addr uint64) Instr {
	return Addu16iD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 16)),
	}
}

func (i Addu16iD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addu16i.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.imm.text())
}

func (i Addu16iD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["addu16i.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.imm.val, 10, 16)

	return writeWord(w, word)
}

func (i Addu16iD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"addu16i.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "imm": i.imm.val},
	)
}
