package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sh - sh rs2, off(rs1).
type Sh struct {
	base

	rs1, rs2 string
	off      imm
}

func (i Sh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sh %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Sh) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["sh"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	return writeWord(w, word)
}

func newSh(ops []Op) (Instr, error) {
	rs2, base, off, err := wantR2M(ops, "sh", false)
	if err != nil {
		return nil, err
	}

	return Sh{
		rs1: base,
		rs2: rs2,
		off: off,
	}, nil
}
