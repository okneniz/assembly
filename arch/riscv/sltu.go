package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sltu - sltu rd, rs1, rs2; pseudo: snez (rs1=zero).
type Sltu struct {
	base

	rd, rs1, rs2 string
}

func decodeSltu(w uint32, addr uint64) Instr {
	return Sltu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i Sltu) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("snez %s, %s", i.rd, i.rs2)
	}

	return fmt.Sprintf("sltu %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sltu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["sltu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20

	return writeWord(w, word)
}

func (i Sltu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sltu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newSltu(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sltu")
	if err != nil {
		return nil, err
	}

	return Sltu{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
