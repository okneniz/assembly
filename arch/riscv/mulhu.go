package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulhu — mulhu rd, rs1, rs2.
type Mulhu struct {
	base

	rd, rs1, rs2 string
}

func decodeMulhu(w uint32, addr uint64) Instr {
	return Mulhu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Mulhu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulhu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Mulhu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["mulhu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mulhu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mulhu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newMulhu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulhu")
	if err != nil {
		return nil, err
	}

	return Mulhu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
