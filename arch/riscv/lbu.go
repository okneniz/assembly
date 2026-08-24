package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lbu — lbu rd, off(rs1).
type Lbu struct {
	base

	rd, rs1 string
	off     imm
}

func decodeLbu(w uint32, addr uint64) Instr {
	return Lbu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Lbu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lbu %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Lbu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lbu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lbu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lbu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
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
