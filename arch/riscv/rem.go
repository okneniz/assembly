package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rem — rem rd, rs1, rs2.
type Rem struct {
	base

	rd, rs1, rs2 string
}

func decodeRem(w uint32, addr uint64) Instr {
	return Rem{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Rem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rem %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Rem) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["rem"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Rem) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"rem",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newRem(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "rem")
	if err != nil {
		return nil, err
	}

	return Rem{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
