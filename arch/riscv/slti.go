package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slti - slti rd, rs1, imm.
type Slti struct {
	rd, rs1 string
	imm     imm
}

func (i Slti) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slti %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Slti) Encode(w io.Writer, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["slti"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|bits)
}

func newSlti(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "slti")
	if err != nil {
		return nil, err
	}

	return Slti{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
