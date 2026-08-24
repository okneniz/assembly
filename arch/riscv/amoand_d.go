package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoandD — amoand.d rd, rs2, (rs1).
type AmoandD struct {
	base

	rd, rs1, rs2 string
}

func decodeAmoandD(w uint32, addr uint64) Instr {
	return AmoandD{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i AmoandD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoand.d %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoandD) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoand_d"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func (i AmoandD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amoand.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64A",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAmoandD(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoand.d: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoand.d: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoand.d: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoand.d")
	if err != nil {
		return nil, err
	}

	return AmoandD{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
