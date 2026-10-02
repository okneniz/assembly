package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srlw - srlw rd, rs1, rs2.
type Srlw struct {
	base

	rd, rs1, rs2 string
}

func (i Srlw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srlw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Srlw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["srlw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSrlw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "srlw")
	if err != nil {
		return nil, err
	}

	return Srlw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
