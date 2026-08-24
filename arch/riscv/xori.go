package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xori - xori rd, rs1, imm; pseudo: not (imm=-1).
type Xori struct {
	base

	rd, rs1 string
	imm     imm
}

func decodeXori(w uint32, addr uint64) Instr {
	return Xori{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func (i Xori) ObjDump(_ disasm.ViewCtx) string {
	if i.imm.val == -1 {
		return fmt.Sprintf("not %s, %s", i.rd, i.rs1)
	}

	return fmt.Sprintf("xori %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Xori) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["xori"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|bits)
}

func (i Xori) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"xori",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newXori(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "xori")
	if err != nil {
		return nil, err
	}

	return Xori{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
