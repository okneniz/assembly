package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Divw - divw rd, rs1, rs2.
type Divw struct {
	base

	rd, rs1, rs2 string
}

// Divw - divw rd, rs1, rs2.
func (Builder) Divw(rd, rs1, rs2 Reg) Instr {
	return Divw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeDivw(w uint32, addr uint64) Instr {
	return Divw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Divw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("divw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Divw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["divw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Divw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"divw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newDivw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "divw")
	if err != nil {
		return nil, err
	}

	return Divw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
