package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sra — sra rd, rs1, rs2.
type Sra struct {
	base

	rd, rs1, rs2 string
}

func decodeSra(w uint32, addr uint64) Instr {
	return Sra{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Sra) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sra %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sra) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["sra"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Sra) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sra",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSra(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sra")
	if err != nil {
		return nil, err
	}

	return Sra{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
