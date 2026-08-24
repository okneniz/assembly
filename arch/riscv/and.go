package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// And — and rd, rs1, rs2.
type And struct {
	base

	rd, rs1, rs2 string
}

func decodeAnd(w uint32, addr uint64) Instr {
	return And{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

// cAnd - compressed forms (c.and): base - halfword, length 2.
func cAnd(h uint32, addr uint64, rd, rs1, rs2 string) And {
	return And{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i And) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("and %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i And) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["and"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 3); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func (i And) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"and",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAnd(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "and")
	if err != nil {
		return nil, err
	}

	return And{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
