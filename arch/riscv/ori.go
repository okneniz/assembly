package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ori - ori rd, rs1, imm.
type Ori struct {
	rd, rs1 string
	imm     imm
}

func (i Ori) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ori %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Ori) Encode(w io.Writer, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["ori"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|bits)
}

func newOri(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "ori")
	if err != nil {
		return nil, err
	}

	return Ori{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
