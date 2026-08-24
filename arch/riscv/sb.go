package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sb — sb rs2, off(rs1).
type Sb struct {
	base

	rs1, rs2 string
	off      imm
}

func decodeSb(w uint32, addr uint64) Instr {
	return Sb{
		base: newBase(addr, w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func (i Sb) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sb %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Sb) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["sb"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	return writeWord(w, word)
}

func (i Sb) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSb(ops []Op) (Instr, error) {
	rs2, base, off, err := wantR2M(ops, "sb", false)
	if err != nil {
		return nil, err
	}

	return Sb{
		rs1: base,
		rs2: rs2,
		off: off,
	}, nil
}
