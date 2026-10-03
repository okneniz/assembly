package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xor - xor rd, rs1, rs2.
type Xor struct {
	rd, rs1, rs2 string
}

// cXor - compressed forms (c.xor): base - halfword, length 2.
func cXor(rd, rs1, rs2 string) Xor {
	return Xor{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}
}

func (i Xor) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["xor"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 1); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func (i Xor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("xor %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newXor(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "xor")
	if err != nil {
		return nil, err
	}

	return Xor{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
