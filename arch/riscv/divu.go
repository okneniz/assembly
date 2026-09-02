package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Divu - divu rd, rs1, rs2.
type Divu struct {
	base

	rd, rs1, rs2 string
}

// Divu - divu rd, rs1, rs2.
func (Builder) Divu(rd, rs1, rs2 Reg) Instr {
	return Divu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeDivu(w uint32, addr uint64) Instr {
	return Divu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Divu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("divu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Divu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["divu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Divu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"divu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newDivu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "divu")
	if err != nil {
		return nil, err
	}

	return Divu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
