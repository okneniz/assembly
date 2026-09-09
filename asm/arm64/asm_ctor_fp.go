package arm64

// Scalar FP assembler constructors: fadd/fsub/fmul/fdiv/fmax/fmin
// (s/d by register type), fneg/fcvt/fcvtzs/scvtf/ucvtf, fmov (registers,
// imm8), fcmp (registers/#0.0), fmadd/fnmsub.

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// newFp3Arm — fop fd, fn, fm: the (d, s) enc pair by the first register's type.
func newFp3Arm(op string, encD, encS uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want fd, fn, fm", op)
		}

		regs, err := armReg3Strings(ops, op)
		if err != nil {
			return nil, err
		}

		fd, fn, fm := regs[0], regs[1], regs[2]
		enc := encS
		if fd[0] == 'd' {
			enc = encD
		}

		return arch.NewFp3(op, fd, fn, fm, enc), nil
	}
}

// armReg3Strings — three FP/integer registers (no strict kind check).
func armReg3Strings(ops []vOp, name string) ([3]string, error) {
	var out [3]string
	if len(ops) != 3 {
		return out, fmt.Errorf("%s: want 3 registers", name)
	}

	for i := range 3 {
		if ops[i].Reg() == "" {
			return out, fmt.Errorf("%s: register operand %d", name, i+1)
		}

		if _, err := armRegNum(ops[i].Reg()); err != nil {
			return out, fmt.Errorf("%s: %w", name, err)
		}

		out[i] = ops[i].Reg()
	}

	return out, nil
}

// newFp2Arm — fop fd, fn (enc by kinds; the rdK/rnK kinds are derived from
// the operands, not from the word).
func newFp2Arm(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want fd, fn", op)
		}

		if ops[0].Reg() == "" || ops[1].Reg() == "" {
			return nil, fmt.Errorf("%s: register operands", op)
		}

		rd, rn := ops[0].Reg(), ops[1].Reg()
		if _, err := armRegNum(rd); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if _, err := armRegNum(rn); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		rdK, rnK := regKindOf(rd), regKindOf(rn)
		return arch.NewFp2(op, rd, rn, enc, rdK, rnK), nil
	}
}

// regKindOf — fpKind by the name prefix.
func regKindOf(name string) fpKind {
	switch name[0] {
	case 's':
		return kS
	case 'd':
		return kD
	case 'w':
		return kW
	default:
		return kX
	}
}

// newFmov — fmov fd, fn (registers) | fmov fd, #imm (imm8 via the
// vfpExpandImm table; text #%.8f).
func newFmov(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("fmov: want fd, op")
	}

	if ops[0].Reg() == "" {
		return nil, errors.New("fmov: register expected")
	}

	rd := ops[0].Reg()
	if _, err := armRegNum(rd); err != nil {
		return nil, fmt.Errorf("fmov: %w", err)
	}

	isS := rd[0] == 's'
	if ops[1].Kind() == arch.ArmOpFloat || ops[1].Kind() == arch.ArmOpImm {
		text := fmt.Sprintf("%.8f", ops[1].Float())
		enc := uint32(0x1E601000)
		rdK := kD
		if isS {
			rdK = kS
		}

		for imm8 := range uint32(256) {
			if isS {
				if fmt.Sprintf("%.8f", vfpExpandImm32(imm8)) == text {
					return arch.NewFmovImm(rd, ops[1].Float(), text, isS, enc, rdK), nil
				}
			} else if fmt.Sprintf("%.8f", vfpExpandImm64(imm8)) == text {
				return arch.NewFmovImm(rd, ops[1].Float(), text, isS, enc, rdK), nil
			}
		}

		return nil, errors.New("fmov: imm not encodable")
	}

	if ops[1].Reg() == "" {
		return nil, errors.New("fmov: register or immediate expected")
	}

	rn := ops[1].Reg()
	if _, err := armRegNum(rn); err != nil {
		return nil, fmt.Errorf("fmov: %w", err)
	}

	// register-move forms: kinds set the encoding (like the decode table)
	rdK, rnK := regKindOf(rd), regKindOf(rn)
	var enc uint32
	switch {
	case rdK == kD && rnK == kD, rdK == kS && rnK == kS:
		enc = 0x1E604000 // fmov d,d / s,s (one base; kinds in the registers)
		if rdK == kS {
			enc = 0x1E204000
		}
	case rdK == kD && rnK == kX:
		enc = 0x9E670000
	case rdK == kX && rnK == kD:
		enc = 0x9E660000
	case rdK == kS && rnK == kW:
		enc = 0x1E270000
	case rdK == kW && rnK == kS:
		enc = 0x1E260000
	default:
		return nil, errors.New("fmov: bad register kinds")
	}

	return arch.NewFp2("fmov", rd, rn, enc, rdK, rnK), nil
}

// newFcmpArm — fcmp fn, fm | fcmp fn, #0.0.
func newFcmpArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("fcmp: want fn, fm|#0.0")
	}

	if ops[0].Reg() == "" {
		return nil, errors.New("fcmp: register expected")
	}

	rn := ops[0].Reg()
	k := regKindOf(rn)
	enc := uint32(0x1E602000)
	if k == kS {
		enc = 0x1E202000
	}

	if ops[1].Kind() == arch.ArmOpFloat || ops[1].Kind() == arch.ArmOpImm {
		return arch.NewFcmp(rn, "", false, enc, enc, k), nil
	}

	if ops[1].Reg() == "" {
		return nil, errors.New("fcmp: register or #0.0")
	}

	return arch.NewFcmp(rn, ops[1].Reg(), true, enc, enc, k), nil
}

// newFmadd — fmadd fd, fn, fm, fa (d form).
func newFmadd(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 4 {
			return nil, fmt.Errorf("%s: want fd, fn, fm, fa", op)
		}

		regs := make([]string, 4)
		for i := range 4 {
			if ops[i].Reg() == "" {
				return nil, fmt.Errorf("%s: register operand %d", op, i+1)
			}

			if _, err := armRegNum(ops[i].Reg()); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

			regs[i] = ops[i].Reg()
		}

		return arch.NewFp4(op, regs[0], regs[1], regs[2], regs[3], enc), nil
	}
}
