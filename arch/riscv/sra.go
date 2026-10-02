package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sra - sra rd, rs1, rs2.
type Sra struct {
	rd, rs1, rs2 string
}

func (i Sra) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sra %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sra) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["sra"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSra(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sra")
	if err != nil {
		return nil, err
	}

	return Sra{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
