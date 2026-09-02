package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Remuw - remuw rd, rs1, rs2.
type Remuw struct {
	base

	rd, rs1, rs2 string
}

// Remuw - remuw rd, rs1, rs2.
func (Builder) Remuw(rd, rs1, rs2 Reg) Instr {
	return Remuw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeRemuw(w uint32, addr uint64) Instr {
	return Remuw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Remuw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("remuw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Remuw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["remuw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Remuw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"remuw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newRemuw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "remuw")
	if err != nil {
		return nil, err
	}

	return Remuw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
