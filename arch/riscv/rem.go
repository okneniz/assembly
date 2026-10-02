package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rem - rem rd, rs1, rs2.
type Rem struct {
	rd, rs1, rs2 string
}

func (i Rem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rem %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Rem) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["rem"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newRem(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "rem")
	if err != nil {
		return nil, err
	}

	return Rem{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
