package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoaddW - amoadd.w rd, rs2, (rs1).
type AmoaddW struct {
	rd, rs1, rs2 string
}

func (i AmoaddW) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoadd_w"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func (i AmoaddW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoadd.w %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func newAmoaddW(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoadd.w: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoadd.w: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoadd.w: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoadd.w")
	if err != nil {
		return nil, err
	}

	return AmoaddW{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
