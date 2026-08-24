package riscv

// Shared operand validators for the new* assembler constructors
// (overlapping instruction families).

import (
	"errors"
	"fmt"
)

// wantR3 - validation of three R-form register operands.
func wantR3(ops []Op, name string) (string, string, string, error) {
	if len(ops) != 3 {
		return "", "", "", fmt.Errorf("%s: want rd, rs1, rs2", name)
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", name, err)
	}

	rs1, err := wantReg(ops[1], false)
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", name, err)
	}

	rs2, err := wantReg(ops[2], false)
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", name, err)
	}

	return rd, rs1, rs2, nil
}

// wantR2T - validation of rs1, rs2 + branch target.
func wantR2T(ops []Op, name string) (string, string, imm, error) {
	if len(ops) != 3 {
		return "", "", imm{}, fmt.Errorf("%s: want rs1, rs2, target", name)
	}

	rs1, err := wantReg(ops[0], false)
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	e, err := wantExpr(ops[2])
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	return rs1, rs2, e, nil
}

// wantR2M - validation of rd/rs2 + mem (load/store).
func wantR2M(ops []Op, name string, fp bool) (string, string, imm, error) {
	if len(ops) != 2 {
		return "", "", imm{}, fmt.Errorf("%s: want 2 operands", name)
	}

	r, err := wantReg(ops[0], fp)
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	base, off, err := wantMem(ops[1], name)
	return r, base, off, err
}

// wantMem - validation of a memory operand: off(reg).
func wantMem(op Op, name string) (string, imm, error) {
	if op.kind != opMemK {
		return "", imm{}, fmt.Errorf("%s: want memory operand", name)
	}

	if _, err := rvRegNumOf(op.mem.base, false); err != nil {
		return "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	return op.mem.base, immNum(op.mem.off), nil
}

// wantI3 - validation of three I-form operands (rd, rs1, imm/shamt).
func wantI3(ops []Op, name string) (string, string, imm, error) {
	if len(ops) != 3 {
		return "", "", imm{}, fmt.Errorf("%s: want rd, rs1, imm", name)
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	rs1, err := wantReg(ops[1], false)
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	e, err := wantExpr(ops[2])
	if err != nil {
		return "", "", imm{}, fmt.Errorf("%s: %w", name, err)
	}

	return rd, rs1, e, nil
}

// wantFP - validation of n floating registers.
func wantFP(ops []Op, name string, n int) ([]string, error) {
	if len(ops) < n {
		return nil, fmt.Errorf("%s: want %d register operands", name, n)
	}

	regs := make([]string, n)
	for k := range n {
		r, err := wantReg(ops[k], true)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		regs[k] = r
	}

	return regs, nil
}

// wantRM - an optional numeric rounding-mode operand (the rne/rtz/...
// keywords are already converted to numbers by the syntax layer).
func wantRM(ops []Op, idx int) imm {
	if idx < len(ops) && ops[idx].kind == opNumK {
		return immNum(ops[idx].num)
	}

	return immNum(0)
}

// wantCSR - a numeric CSR operand (the name is already converted by the
// syntax layer; the range is checked here - computed instruction construction).
func wantCSR(op Op) (int64, error) {
	if op.kind != opNumK {
		return 0, errors.New("CSR operand expected")
	}

	if op.num < 0 || op.num > 0xfff {
		return 0, fmt.Errorf("CSR address %#x out of 12-bit range", op.num)
	}

	return op.num, nil
}
