package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Divw - divw rd, rs1, rs2.
type Divw struct {
	rd, rs1, rs2 string
}

func (i Divw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["divw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Divw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("divw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newDivw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "divw")
	if err != nil {
		return nil, err
	}

	return Divw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
