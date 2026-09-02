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

// Sh - sh rs2, off(rs1).
func (Builder) Sh(rs2, rs1 Reg, off Off) Instr {
	return Sh{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

func decodeSh(w uint32, addr uint64) Instr {
	return Sh{
		base: newBase(addr, w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func (i Sh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sh %s, %s(%s)", i.rs2, i.off.text(), i.rs1)
}

func (i Sh) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encS(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["sh"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	return writeWord(w, word)
}

func (i Sh) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
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
