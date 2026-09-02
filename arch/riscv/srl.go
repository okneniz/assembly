package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srl - srl rd, rs1, rs2.
type Srl struct {
	base

	rd, rs1, rs2 string
}

// Srl - srl rd, rs1, rs2.
func (Builder) Srl(rd, rs1, rs2 Reg) Instr {
	return Srl{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeSrl(w uint32, addr uint64) Instr {
	return Srl{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Srl) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srl %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Srl) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["srl"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Srl) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srl",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSrl(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "srl")
	if err != nil {
		return nil, err
	}

	return Srl{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
