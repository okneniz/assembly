package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lb - lb rd, off(rs1).
type Lb struct {
	rd, rs1 string
	off     imm
}

func (i Lb) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lb"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lb) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lb %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func newLb(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lb", false)
	if err != nil {
		return nil, err
	}

	return Lb{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
