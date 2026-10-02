package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sllw - sllw rd, rs1, rs2.
type Sllw struct {
	base

	rd, rs1, rs2 string
}

func (i Sllw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sllw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sllw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["sllw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSllw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sllw")
	if err != nil {
		return nil, err
	}

	return Sllw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
