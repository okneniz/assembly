package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulhu - mulhu rd, rs1, rs2.
type Mulhu struct {
	rd, rs1, rs2 string
}

func (i Mulhu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulhu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Mulhu) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["mulhu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newMulhu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulhu")
	if err != nil {
		return nil, err
	}

	return Mulhu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
