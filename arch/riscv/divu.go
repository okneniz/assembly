package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Divu - divu rd, rs1, rs2.
type Divu struct {
	rd, rs1, rs2 string
}

func (i Divu) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["divu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Divu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("divu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newDivu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "divu")
	if err != nil {
		return nil, err
	}

	return Divu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
