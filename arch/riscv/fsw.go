package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fsw - fsw rs2, off(rs1) (rs2 - floating).
type Fsw struct {
	base

	rs1, rs2 string
	off      imm
}

func decodeFsw(w uint32) Instr {
	return Fsw{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func (i Fsw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fsw %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Fsw) Encode(w io.Writer, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["fsw"][0] | regBits(i.rs1)<<15 | fregBits(i.rs2)<<20 | bits
	return writeWord(w, word)
}

func newFsw(ops []Op) (Instr, error) {
	rs2, base, off, err := wantR2M(ops, "fsw", true)
	if err != nil {
		return nil, err
	}

	return Fsw{
		rs1: base,
		rs2: rs2,
		off: off,
	}, nil
}
