package arm64

// Common operand validators and small helpers of the new* assembler
// constructors (overlapping instruction families).

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

// armRegNum parses a register name: prefix + number (x0..x30, w3, v31, d7,
// s2, b0, h5, q1). Named ones (sp/xzr/wzr/wsp) → 31.
func armRegNum(name string) (uint32, error) {
	switch name {
	case "sp", "xzr", "wzr", "wsp":
		return 31, nil
	}

	if len(name) < 2 {
		return 0, fmt.Errorf("bad register %q", name)
	}

	n, err := strconv.Atoi(name[1:])
	if err != nil || n < 0 || n > 31 {
		return 0, fmt.Errorf("bad register %q", name)
	}

	return uint32(n), nil
}

// shiftAmt — the shift amount of a modifier operand.
func shiftAmt(op vOp) int64 { return op.num }

// invSysReg: system register name → 15-bit key (inverse of sysregNames;
// the objdump form S<op0>_<op1>_C<n>_C<m>_<op2> is also accepted).
var invSysRegNames = func() map[string]uint32 {
	m := map[string]uint32{}
	for k, name := range sysregNames {
		if _, exists := m[name]; !exists {
			m[name] = k
		}
	}

	return m
}()

func invSysReg(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("sysreg expected")
	}

	if k, ok := invSysRegNames[s]; ok {
		return k, nil
	}

	if strings.HasPrefix(s, "S") && strings.Count(s, "_") == 4 {
		parts := strings.Split(s[1:], "_")
		op0, err0 := strconv.Atoi(parts[0])
		op1, err1 := strconv.Atoi(parts[1])
		crn, err2 := strconv.Atoi(strings.TrimPrefix(parts[2], "C"))
		crm, err3 := strconv.Atoi(strings.TrimPrefix(parts[3], "C"))
		op2, err4 := strconv.Atoi(parts[4])
		if err0 == nil && err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
			op0 >= 2 && op0 <= 3 && crn >= 0 && crn <= 15 && crm >= 0 && crm <= 15 {
			return uint32(
				op0&1,
			)<<14 | uint32(
				op1,
			)<<11 | uint32(
				crn,
			)<<7 | uint32(
				crm,
			)<<3 | uint32(
				op2,
			), nil
		}
	}

	return 0, fmt.Errorf("unknown system register %q", s)
}

// arrQSize — (Q, size) from the .8b/.16b/.4h/.8h/.2s/.4s/.2d suffix.
func arrQSize(arr string) (q, size uint32, err error) {
	switch arr {
	case "8b":
		return 0, 0, nil
	case "16b":
		return 1, 0, nil
	case "4h":
		return 0, 1, nil
	case "8h":
		return 1, 1, nil
	case "2s":
		return 0, 2, nil
	case "4s":
		return 1, 2, nil
	case "2d":
		return 1, 3, nil
	}

	return 0, 0, fmt.Errorf("unknown arrangement %q", arr)
}

// WantTarget — a computed target operand, unwrapped to int64.
func WantTarget(op ArmOp, name string) (int64, error) {
	v, err := wantTarget(op, name)
	if err != nil {
		return 0, err
	}

	return v.val, nil
}
