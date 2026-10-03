package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lwu - lwu rd, off(rs1).
type Lwu struct {
	rd, rs1 string
	off     imm
}

func (i Lwu) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lwu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lwu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lwu %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func newLwu(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lwu", false)
	if err != nil {
		return nil, err
	}

	return Lwu{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
