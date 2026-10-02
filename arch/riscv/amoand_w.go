package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoandW - amoand.w rd, rs2, (rs1).
type AmoandW struct {
	rd, rs1, rs2 string
}

func (i AmoandW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoand.w %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoandW) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoand_w"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
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
