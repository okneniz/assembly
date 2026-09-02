package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lh - lh rd, off(rs1).
type Lh struct {
	base

	rd, rs1 string
	off     imm
}

// Lh - lh rd, off(rs1).
func (Builder) Lh(rd, rs1 Reg, off Off) Instr {
	return Lh{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeLh(w uint32, addr uint64) Instr {
	return Lh{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Lh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lh %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Lh) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lh"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lh) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newLh(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lh", false)
	if err != nil {
		return nil, err
	}

	return Lh{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
