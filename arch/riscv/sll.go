package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sll - sll rd, rs1, rs2.
type Sll struct {
	rd, rs1, rs2 string
}

func (i Sll) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sll %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sll) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["sll"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSll(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sll")
	if err != nil {
		return nil, err
	}

	return Sll{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
