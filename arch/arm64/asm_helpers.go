package arm64

// Common operand validators and small helpers of the new* assembler
// constructors (overlapping instruction families).

import "fmt"

// wantTarget — a target operand (a computed value).
func wantTarget(op vOp, name string) (imm, error) {
	if op.kind != armOpImm || op.sym != "" {
		return imm{}, fmt.Errorf("%s: want target operand", name)
	}

	return immNum(op.num), nil
}

// wantAReg — an x/w register operand.
func wantAReg(op vOp, name string) (string, error) {
	if op.kind != armOpReg || op.reg == "" {
		return "", fmt.Errorf("%s: want register operand", name)
	}

	if _, err := armRegNum(op.reg); err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}

	return op.reg, nil
}

// armReg2/3 — validation of register operands from parsing.
func armReg2(ops []vOp, name string) (string, string, error) {
	r0, err := wantAReg(ops[0], name)
	if err != nil {
		return "", "", err
	}

	r1, err := wantAReg(ops[1], name)
	if err != nil {
		return "", "", err
	}

	return r0, r1, nil
}

func armReg3(ops []vOp, name string) (string, string, string, error) {
	r0, r1, err := armReg2(ops, name)
	if err != nil {
		return "", "", "", err
	}

	r2, err := wantAReg(ops[2], name)
	if err != nil {
		return "", "", "", err
	}

	return r0, r1, r2, nil
}

// wantCond — a condition operand (a name operand).
func wantCond(op vOp, name string) (string, error) {
	if op.sym == "" {
		return "", fmt.Errorf("%s: condition expected", name)
	}

	if _, err := condNum(op.sym); err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}

	return op.sym, nil
}

// zeroReg — xzr/wzr by width.
func zeroReg(rd string) string {
	if rd[0] == 'w' {
		return "wzr"
	}

	return "xzr"
}
