package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Remu - remu rd, rs1, rs2.
type Remu struct {
	base

	rd, rs1, rs2 string
}

// Remu - remu rd, rs1, rs2.
func (Builder) Remu(rd, rs1, rs2 Reg) Instr {
	return Remu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeRemu(w uint32, addr uint64) Instr {
	return Remu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Remu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("remu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Remu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["remu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Remu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"remu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newRemu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "remu")
	if err != nil {
		return nil, err
	}

	return Remu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
