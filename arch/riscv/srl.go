package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srl - srl rd, rs1, rs2.
type Srl struct {
	base

	rd, rs1, rs2 string
}

func (i Srl) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srl %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Srl) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["srl"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSrl(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "srl")
	if err != nil {
		return nil, err
	}

	return Srl{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
