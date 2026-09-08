package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoswapW - amoswap.w rd, rs2, (rs1).
type AmoswapW struct {
	base

	rd, rs1, rs2 string
}

// AmoswapW - amoswap.w rd, rs2, (rs1).
func (Builder) AmoswapW(rd, rs1, rs2 Reg) Instr {
	return AmoswapW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeAmoswapW(w uint32) Instr {
	return AmoswapW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i AmoswapW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoswap.w %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoswapW) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoswap_w"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func newAmoswapW(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoswap.w: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoswap.w: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoswap.w: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoswap.w")
	if err != nil {
		return nil, err
	}

	return AmoswapW{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
