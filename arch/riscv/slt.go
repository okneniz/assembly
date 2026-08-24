package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slt — slt rd, rs1, rs2.
type Slt struct {
	base

	rd, rs1, rs2 string
}

func decodeSlt(w uint32, addr uint64) Instr {
	return Slt{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Slt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slt %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Slt) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["slt"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Slt) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"slt",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSlt(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "slt")
	if err != nil {
		return nil, err
	}

	return Slt{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
