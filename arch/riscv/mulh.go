package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulh - mulh rd, rs1, rs2.
type Mulh struct {
	rd, rs1, rs2 string
}

func (i Mulh) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["mulh"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newMulh(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulh")
	if err != nil {
		return nil, err
	}

	return Mulh{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
