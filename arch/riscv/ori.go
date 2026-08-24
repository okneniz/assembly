package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ori — ori rd, rs1, imm.
type Ori struct {
	base

	rd, rs1 string
	imm     imm
}

func decodeOri(w uint32, addr uint64) Instr {
	return Ori{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func (i Ori) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ori %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Ori) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["ori"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|bits)
}

func (i Ori) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ori",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newOri(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "ori")
	if err != nil {
		return nil, err
	}

	return Ori{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
