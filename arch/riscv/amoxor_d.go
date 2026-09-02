package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoxorD - amoxor.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] ^= rs2.
type AmoxorD struct {
	base

	rd, rs1, rs2 string
}

// AmoxorD - amoxor.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] ^= rs2.
func (Builder) AmoxorD(rd, rs1, rs2 Reg) Instr {
	return AmoxorD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeAmoxorD(w uint32, addr uint64) Instr {
	return AmoxorD{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i AmoxorD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoxor.d %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoxorD) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoxor_d"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func (i AmoxorD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amoxor.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64A",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAmoxorD(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoxor.d: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoxor.d: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoxor.d: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoxor.d")
	if err != nil {
		return nil, err
	}

	return AmoxorD{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
