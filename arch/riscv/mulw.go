package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulw - mulw rd, rs1, rs2.
type Mulw struct {
	rd, rs1, rs2 string
}

func (i Mulw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["mulw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mulw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newMulw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulw")
	if err != nil {
		return nil, err
	}

	return Mulw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
