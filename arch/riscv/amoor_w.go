package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoorW - amoor.w rd, rs2, (rs1).
type AmoorW struct {
	base

	rd, rs1, rs2 string
}

func (i AmoorW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoor.w %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoorW) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoor_w"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func newAmoorW(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoor.w: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoor.w: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoor.w: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoor.w")
	if err != nil {
		return nil, err
	}

	return AmoorW{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
