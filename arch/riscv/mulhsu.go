package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulhsu - mulhsu rd, rs1, rs2.
type Mulhsu struct {
	base

	rd, rs1, rs2 string
}

func (i Mulhsu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulhsu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Mulhsu) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["mulhsu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newMulhsu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulhsu")
	if err != nil {
		return nil, err
	}

	return Mulhsu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
