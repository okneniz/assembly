package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lbu - lbu rd, off(rs1).
type Lbu struct {
	rd, rs1 string
	off     imm
}

func (i Lbu) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lbu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lbu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lbu %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func newLbu(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lbu", false)
	if err != nil {
		return nil, err
	}

	return Lbu{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
