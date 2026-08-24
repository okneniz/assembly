package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sllw — sllw rd, rs1, rs2.
type Sllw struct {
	base

	rd, rs1, rs2 string
}

func decodeSllw(w uint32, addr uint64) Instr {
	return Sllw{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Sllw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sllw %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sllw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["sllw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Sllw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sllw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSllw(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sllw")
	if err != nil {
		return nil, err
	}

	return Sllw{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
