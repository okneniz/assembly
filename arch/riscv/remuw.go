package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Remuw - remuw rd, rs1, rs2.
type Remuw struct {
	rd, rs1, rs2 string
}

func (i Remuw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("remuw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Remuw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["remuw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newRemuw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "remuw")
	if err != nil {
		return nil, err
	}

	return Remuw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
