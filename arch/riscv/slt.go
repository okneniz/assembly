package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slt - slt rd, rs1, rs2; pseudo: sltz (rs2 = x0), sgtz (rs1 = x0).
type Slt struct {
	rd, rs1, rs2 string
}

func (i Slt) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["slt"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Slt) ObjDump(_ disasm.ViewCtx) string {
	switch {
	case i.rs2 == "zero":
		return fmt.Sprintf("sltz %s, %s", i.rd, i.rs1)
	case i.rs1 == "zero":
		return fmt.Sprintf("sgtz %s, %s", i.rd, i.rs2)
	}

	return fmt.Sprintf("slt %s, %s, %s", i.rd, i.rs1, i.rs2)
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
