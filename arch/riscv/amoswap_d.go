package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoswapD - amoswap.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] = rs2.
type AmoswapD struct {
	rd, rs1, rs2 string
}

func (i AmoswapD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoswap.d %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoswapD) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoswap_d"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func newAmoswapD(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoswap.d: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoswap.d: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoswap.d: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoswap.d")
	if err != nil {
		return nil, err
	}

	return AmoswapD{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
