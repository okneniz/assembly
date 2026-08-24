package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lhu — lhu rd, off(rs1).
type Lhu struct {
	base

	rd, rs1 string
	off     imm
}

func decodeLhu(w uint32, addr uint64) Instr {
	return Lhu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Lhu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lhu %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Lhu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lhu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lhu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lhu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newLhu(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lhu", false)
	if err != nil {
		return nil, err
	}

	return Lhu{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
