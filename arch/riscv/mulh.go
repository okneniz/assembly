package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulh - mulh rd, rs1, rs2.
type Mulh struct {
	base

	rd, rs1, rs2 string
}

// Mulh - mulh rd, rs1, rs2.
func (Builder) Mulh(rd, rs1, rs2 Reg) Instr {
	return Mulh{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeMulh(w uint32, addr uint64) Instr {
	return Mulh{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Mulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulh %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Mulh) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["mulh"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mulh) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mulh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newMulh(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulh")
	if err != nil {
		return nil, err
	}

	return Mulh{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
