package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Remu - remu rd, rs1, rs2.
type Remu struct {
	rd, rs1, rs2 string
}

func (i Remu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("remu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Remu) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["remu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newRemu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "remu")
	if err != nil {
		return nil, err
	}

	return Remu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
