package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mulw - mulw rd, rs1, rs2.
type Mulw struct {
	base

	rd, rs1, rs2 string
}

// Mulw - mulw rd, rs1, rs2.
func (Builder) Mulw(rd, rs1, rs2 Reg) Instr {
	return Mulw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeMulw(w uint32, addr uint64) Instr {
	return Mulw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Mulw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mulw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Mulw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["mulw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Mulw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mulw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newMulw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "mulw")
	if err != nil {
		return nil, err
	}

	return Mulw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
