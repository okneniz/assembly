package arm64

// Scalar FP assembler constructors: fadd/fsub/fmul/fdiv/fmax/fmin
// (s/d by register type), fneg/fcvt, the int↔FP conversions scvtf/
// ucvtf/fcvtzs/fcvtzu, fmov (registers, imm8, GPR moves), fcmp
// (registers/#0.0), fmadd/fnmsub. Every entry builds through the arch
// Builder — the one construction API of the package.

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// fpReg — one FP register operand (s/d by the spelling).
func fpReg(op vOp, name string) (arch.FReg, error) {
	r, err := arch.FRegOf(op.Reg())
	if err != nil {
		return arch.FReg{}, fmt.Errorf("%s: %w", name, err)
	}

	return r, nil
}

// gpr — one integer register operand (w/x, zr allowed).
func gpr(op vOp, name string) (arch.Reg, error) {
	if op.Reg() == "" {
		return arch.Reg{}, fmt.Errorf("%s: register operand expected", name)
	}

	r, err := arch.RegOf(op.Reg())
	if err != nil {
		return arch.Reg{}, fmt.Errorf("%s: %w", name, err)
	}

	return r, nil
}

// newFp3Arm — an FP three-register instruction (fd, fn, fm): the
// Builder method comes in as a method expression.
func newFp3Arm(
	name string,
	method func(arch.Builder, arch.FReg, arch.FReg, arch.FReg) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want fd, fn, fm", name)
		}

		rd, err := fpReg(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := fpReg(ops[1], name)
		if err != nil {
			return nil, err
		}

		rm, err := fpReg(ops[2], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn, rm)
	}
}

// newFp2Arm — an FP two-register instruction (fd, fn).
func newFp2Arm(
	name string,
	method func(arch.Builder, arch.FReg, arch.FReg) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want fd, fn", name)
		}

		rd, err := fpReg(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := fpReg(ops[1], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn)
	}
}

// newFpConvArm — an int→FP conversion (fd, wn|xn).
func newFpConvArm(
	name string,
	method func(arch.Builder, arch.FReg, arch.Reg) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want fd, wn|xn", name)
		}

		rd, err := fpReg(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := gpr(ops[1], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn)
	}
}

// newConvFpArm — an FP→int conversion (wd|xd, fn).
func newConvFpArm(
	name string,
	method func(arch.Builder, arch.Reg, arch.FReg) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want wd|xd, fn", name)
		}

		rd, err := gpr(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := fpReg(ops[1], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn)
	}
}

// newFp4Arm — an FP four-register instruction (fd, fn, fm, fa).
func newFp4Arm(
	name string,
	method func(arch.Builder, arch.FReg, arch.FReg, arch.FReg, arch.FReg) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 4 {
			return nil, fmt.Errorf("%s: want fd, fn, fm, fa", name)
		}

		rd, err := fpReg(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := fpReg(ops[1], name)
		if err != nil {
			return nil, err
		}

		rm, err := fpReg(ops[2], name)
		if err != nil {
			return nil, err
		}

		ra, err := fpReg(ops[3], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn, rm, ra)
	}
}

// newFmov — fmov by the operand spellings: fd, fn (two FP registers),
// fd, xn | sn, wn (from a GPR), xn, fd | wn, sn (to a GPR), fd, #imm
// (the imm8 form).
func newFmov(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("fmov: want fd, op")
	}

	rd, rdIsFp := arch.FRegOf(ops[0].Reg())

	if ops[1].Kind() == arch.ArmOpFloat || ops[1].Kind() == arch.ArmOpImm {
		if ops[0].Reg() == "" || rdIsFp != nil {
			return nil, errors.New("fmov: FP register expected")
		}

		return (arch.Builder{}).FmovImm(rd, ops[1].Float())
	}

	rn, rnIsFp := arch.FRegOf(ops[1].Reg())
	switch {
	case rdIsFp == nil && rnIsFp == nil:
		return (arch.Builder{}).Fmov(rd, rn)
	case rdIsFp == nil:
		src, err := gpr(ops[1], "fmov")
		if err != nil {
			return nil, err
		}

		return (arch.Builder{}).FmovFromGpr(rd, src)
	case rnIsFp == nil:
		dst, err := gpr(ops[0], "fmov")
		if err != nil {
			return nil, err
		}

		return (arch.Builder{}).FmovToGpr(dst, rn)
	}

	return nil, errors.New("fmov: register or immediate expected")
}

// newFcmpArm — fcmp fn, fm | fcmp fn, #0.0.
func newFcmpArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("fcmp: want fn, fm|#0.0")
	}

	rn, err := fpReg(ops[0], "fcmp")
	if err != nil {
		return nil, err
	}

	if ops[1].Kind() == arch.ArmOpFloat || ops[1].Kind() == arch.ArmOpImm {
		return (arch.Builder{}).FcmpZero(rn)
	}

	rm, err := fpReg(ops[1], "fcmp")
	if err != nil {
		return nil, errors.New("fcmp: FP register or #0.0")
	}

	return (arch.Builder{}).Fcmp(rn, rm)
}
