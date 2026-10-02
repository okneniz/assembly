package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Addw - addw rd, rs1, rs2.
type Addw struct {
	base

	rd, rs1, rs2 string
}

// cAddw - compressed forms (c.addw): base - halfword, length 2.
func cAddw(h uint32, rd, rs1, rs2 string) Addw {
	return Addw{
		base: newHalfBase(h),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Addw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Addw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["addw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x9C01, 1); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func newAddw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "addw")
	if err != nil {
		return nil, err
	}

	return Addw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
