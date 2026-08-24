package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Or — or rd, rs1, rs2.
type Or struct {
	base

	rd, rs1, rs2 string
}

func decodeOr(w uint32, addr uint64) Instr {
	return Or{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

// cOr - compressed forms (c.or): base - halfword, length 2.
func cOr(h uint32, addr uint64, rd, rs1, rs2 string) Or {
	return Or{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Or) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("or %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Or) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["or"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 2); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func (i Or) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"or",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newOr(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "or")
	if err != nil {
		return nil, err
	}

	return Or{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
