package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Xor - xor rd, rs1, rs2.
type Xor struct {
	base

	rd, rs1, rs2 string
}

// Xor - xor rd, rs1, rs2.
func (Builder) Xor(rd, rs1, rs2 Reg) Instr {
	return Xor{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeXor(w uint32, addr uint64) Instr {
	return Xor{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

// cXor - compressed forms (c.xor): base - halfword, length 2.
func cXor(h uint32, addr uint64, rd, rs1, rs2 string) Xor {
	return Xor{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Xor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("xor %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Xor) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["xor"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 1); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func (i Xor) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"xor",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newXor(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "xor")
	if err != nil {
		return nil, err
	}

	return Xor{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
