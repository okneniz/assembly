package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sraw - sraw rd, rs1, rs2.
type Sraw struct {
	rd, rs1, rs2 string
}

func (i Sraw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sraw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sraw) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["sraw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func newSraw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sraw")
	if err != nil {
		return nil, err
	}

	return Sraw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
