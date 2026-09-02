package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoandW - amoand.w rd, rs2, (rs1).
type AmoandW struct {
	base

	rd, rs1, rs2 string
}

// AmoandW - amoand.w rd, rs2, (rs1).
func (Builder) AmoandW(rd, rs1, rs2 Reg) Instr {
	return AmoandW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeAmoandW(w uint32, addr uint64) Instr {
	return AmoandW{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i AmoandW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoand.w %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoandW) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoand_w"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func (i AmoandW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amoand.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32A",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAmoandW(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoand.w: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoand.w: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoand.w: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoand.w")
	if err != nil {
		return nil, err
	}

	return AmoandW{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
