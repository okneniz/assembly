package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mul - mul rd, rs1, rs2.
type Mul struct {
	rd, rs1, rs2 string
}

func (i Mul) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["mul"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mul) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mul %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newMul(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mul")
	if err != nil {
		return nil, err
	}

	return Mul{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
