package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Divuw - divuw rd, rs1, rs2.
type Divuw struct {
	base

	rd, rs1, rs2 string
}

// Divuw - divuw rd, rs1, rs2.
func (Builder) Divuw(rd, rs1, rs2 Reg) Instr {
	return Divuw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeDivuw(w uint32, addr uint64) Instr {
	return Divuw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Divuw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("divuw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Divuw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["divuw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Divuw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"divuw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newDivuw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "divuw")
	if err != nil {
		return nil, err
	}

	return Divuw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
