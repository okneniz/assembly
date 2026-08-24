package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fsd - fsd rs2, off(rs1) (rs2 - floating).
type Fsd struct {
	base

	rs1, rs2 string
	off      imm
}

func decodeFsd(w uint32, addr uint64) Instr {
	return Fsd{
		base: newBase(addr, w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func (i Fsd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fsd %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Fsd) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["fsd"][0] | regBits(i.rs1)<<15 | fregBits(i.rs2)<<20 | bits
	return writeWord(w, word)
}

func (i Fsd) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"fsd",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64D",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newFsd(ops []Op) (Instr, error) {
	rs2, base, off, err := wantR2M(ops, "fsd", true)
	if err != nil {
		return nil, err
	}

	return Fsd{
		rs1: base,
		rs2: rs2,
		off: off,
	}, nil
}
