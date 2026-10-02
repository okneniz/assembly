package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slt - slt rd, rs1, rs2.
type Slt struct {
	base

	rd, rs1, rs2 string
}

func (i Slt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slt %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Slt) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["slt"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSlt(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "slt")
	if err != nil {
		return nil, err
	}

	return Slt{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
