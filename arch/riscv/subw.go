package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Subw - subw rd, rs1, rs2; pseudo: negw (rs1=zero).
type Subw struct {
	base

	rd, rs1, rs2 string
}

// cSubw - compressed forms (c.subw): base - halfword, length 2.
func cSubw(h uint32, rd, rs1, rs2 string) Subw {
	return Subw{
		base: newHalfBase(h),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Subw) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("negw %s, %s", i.rd, i.rs2)
	}

	return fmt.Sprintf("subw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Subw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["subw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x9C01, 0); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func newSubw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "subw")
	if err != nil {
		return nil, err
	}

	return Subw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
