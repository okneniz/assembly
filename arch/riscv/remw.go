package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Remw - remw rd, rs1, rs2.
type Remw struct {
	base

	rd, rs1, rs2 string
}

func (i Remw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("remw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Remw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["remw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newRemw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "remw")
	if err != nil {
		return nil, err
	}

	return Remw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
