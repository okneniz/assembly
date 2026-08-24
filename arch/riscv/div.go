package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Div — div rd, rs1, rs2.
type Div struct {
	base

	rd, rs1, rs2 string
}

func decodeDiv(w uint32, addr uint64) Instr {
	return Div{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Div) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("div %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Div) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["div"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Div) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"div",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32M",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newDiv(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "div")
	if err != nil {
		return nil, err
	}

	return Div{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
