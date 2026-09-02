package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Subw - subw rd, rs1, rs2; pseudo: negw (rs1=zero).
type Subw struct {
	base

	rd, rs1, rs2 string
}

// Subw - subw rd, rs1, rs2 (rs1 = zero is printed as negw).
func (Builder) Subw(rd, rs1, rs2 Reg) Instr {
	return Subw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeSubw(w uint32, addr uint64) Instr {
	return Subw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

// cSubw - compressed forms (c.subw): base - halfword, length 2.
func cSubw(h uint32, addr uint64, rd, rs1, rs2 string) Subw {
	return Subw{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Subw) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("negw %s, %s", i.rd, i.rs2)
	}

	return fmt.Sprintf("subw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Subw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["subw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x9C01, 0); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func (i Subw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"subw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSubw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "subw")
	if err != nil {
		return nil, err
	}

	return Subw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
