package arm64

// Logical-family assembler constructors: and/ands/orr/eor/bic/bics/orn/
// eon (bitmask immediate via encodeBitMasks and register shifts). The
// aliases tst (ands Rd=zr), mvn (orn Rn=zr), mov (orr Rn=zr; immediate —
// movz/movn by the legacy hw layout) — asm/arm64/alias. Vector operands
// and symbolic immediates — a constructor error (falling through to the
// legacy candidates).

import (
	"fmt"
)

// isGPR — an integer register (the x/w family, sp/zr).
func isGPR(name string) bool {
	return name != "" && (name[0] == 'x' || name[0] == 'w')
}

// logCtor — rd, rn, X{, lsl #N} / #bitmask.
func logCtor(ops []vOp, name string) (Instr, error) {
	if len(ops) != 3 && len(ops) != 4 {
		return nil, fmt.Errorf("%s: want rd, rn, op", name)
	}

	rd, rn, err := armReg2(ops, name)
	if err != nil {
		return nil, err
	}

	if !isGPR(rd) || !isGPR(rn) {
		return nil, fmt.Errorf("%s: integer registers expected", name)
	}

	op := ops[2]
	if op.kind == armOpImm {
		if len(ops) == 4 {
			return nil, fmt.Errorf("%s: no modifier for immediate form", name)
		}

		if op.kind != armOpImm || op.sym != "" {
			return nil, fmt.Errorf("%s: immediate expected", name)
		}

		v := op.num

		is64 := rd[0] == 'x'
		n, immr, imms, ok := encodeBitMasks(is64, uint64(v))
		if !ok {
			return nil, fmt.Errorf("%s: %#x not encodable as bitmask", name, v)
		}

		return makeLogImmStruct(name, rd, rn, immr, imms, n == 1, is64), nil
	}

	if op.kind == armOpReg && isGPR(op.reg) {
		rm := op.reg
		shift, amt := "lsl", uint32(0)
		if len(ops) == 4 {
			if ops[3].kind != armOpShift {
				return nil, fmt.Errorf("%s: shift modifier expected", name)
			}

			a := shiftAmt(ops[3])
			if a < 0 || a > 63 {
				return nil, fmt.Errorf("%s: bad shift", name)
			}

			shift, amt = ops[3].shift, uint32(a)
		}

		return makeLogShiftStruct(name, rd, rn, rm, shift, amt, rd[0] == 'x'), nil
	}

	return nil, fmt.Errorf("%s: bad operand", name)
}

// makeLogImmStruct/makeLogShiftStruct — assembly by the base name.
func makeLogImmStruct(name, rd, rn string, immr, imms uint32, n, is64 bool) Instr {
	base := newLogImm(rd, rn, immr, imms, n, is64)
	switch name {
	case "orr":
		return OrrImm{logImm: base}
	case "eor":
		return EorImm{logImm: base}
	case "ands":
		return AndsImm{logImm: base}
	default:
		return AndImm{logImm: base}
	}
}

func makeLogShiftStruct(name, rd, rn, rm, shift string, amt uint32, isf bool) Instr {
	common := struct {
		rd, rn, rm string
		imm6       uint32
		shift      string
		isf        bool
	}{
		rd,
		rn,
		rm,
		amt,
		shift,
		isf,
	}
	switch name {
	case "ands":
		return AndsShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	case "bic":
		return BicShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	case "orr":
		return OrrShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	case "orn":
		return OrnShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	case "eor":
		return EorShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	case "eon":
		return EonShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	case "bics":
		return BicsShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	default:
		return AndShift{
			rd:    common.rd,
			rn:    common.rn,
			rm:    common.rm,
			imm6:  common.imm6,
			shift: common.shift,
			isf:   common.isf,
		}
	}
}

func newAndArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "and")
}
func newAndsArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "ands")
}
func newOrrArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "orr")
}
func newEorArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "eor")
}
func newBicArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "bic")
}
func newBicsArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "bics")
}
func newOrnArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "orn")
}
func newEonArm(ops []vOp) (Instr, error) {
	return logCtor(ops, "eon")
}
