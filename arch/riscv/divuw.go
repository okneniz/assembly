package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Divuw - divuw rd, rs1, rs2.
type Divuw struct {
	rd, rs1, rs2 string
}

func (i Divuw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["divuw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Divuw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("divuw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newDivuw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "divuw")
	if err != nil {
		return nil, err
	}

	return Divuw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
