package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sraw - sraw rd, rs1, rs2.
type Sraw struct {
	base

	rd, rs1, rs2 string
}

// Sraw - sraw rd, rs1, rs2.
func (Builder) Sraw(rd, rs1, rs2 Reg) Instr {
	return Sraw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeSraw(w uint32, addr uint64) Instr {
	return Sraw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Sraw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sraw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sraw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["sraw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Sraw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sraw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSraw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sraw")
	if err != nil {
		return nil, err
	}

	return Sraw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
