package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sub - sub rd, rs1, rs2; pseudo: neg (rs1=zero).
type Sub struct {
	base

	rd, rs1, rs2 string
}

// Sub - sub rd, rs1, rs2 (rs1 = zero is printed as neg).
func (Builder) Sub(rd, rs1, rs2 Reg) Instr {
	return Sub{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeSub(w uint32, addr uint64) Instr {
	return Sub{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

// cSub - compressed forms (c.sub): base - halfword, length 2.
func cSub(h uint32, addr uint64, rd, rs1, rs2 string) Sub {
	return Sub{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Sub) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("neg %s, %s", i.rd, i.rs2)
	}

	return fmt.Sprintf("sub %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sub) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["sub"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 0); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func (i Sub) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sub",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSub(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sub")
	if err != nil {
		return nil, err
	}

	return Sub{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
