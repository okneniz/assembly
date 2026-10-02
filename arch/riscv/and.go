package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// And - and rd, rs1, rs2.
type And struct {
	rd, rs1, rs2 string
}

// cAnd - compressed forms (c.and): base - halfword, length 2.
func cAnd(rd, rs1, rs2 string) And {
	return And{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}
}

func (i And) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("and %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i And) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["and"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 3); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func newAnd(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "and")
	if err != nil {
		return nil, err
	}

	return And{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
