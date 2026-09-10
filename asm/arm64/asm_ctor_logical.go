package arm64

// Logical-family assembler constructors: and/ands/orr/eor/bic/bics/orn/
// eon (bitmask immediate via encodeBitMasks and register shifts). The
// aliases tst (ands Rd=zr), mvn (orn Rn=zr), mov (orr Rn=zr; immediate —
// movz/movn by the legacy hw layout) — asm/arm64/alias. Vector operands
// and symbolic immediates — a constructor error (falling through to the
// legacy candidates).

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
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
	if op.Kind() == arch.ArmOpImm {
		if len(ops) == 4 {
			return nil, fmt.Errorf("%s: no modifier for immediate form", name)
		}

		if op.Kind() != arch.ArmOpImm || op.Sym() != "" {
			return nil, fmt.Errorf("%s: immediate expected", name)
		}

		v := op.Num()

		is64 := rd[0] == 'x'
		n, immr, imms, ok := encodeBitMasks(is64, uint64(v))
		if !ok {
			return nil, fmt.Errorf("%s: %#x not encodable as bitmask", name, v)
		}

		return makeLogImmStruct(name, rd, rn, immr, imms, n == 1, is64)
	}

	if op.Kind() == arch.ArmOpReg && isGPR(op.Reg()) {
		rm := op.Reg()
		shift, amt := "lsl", uint32(0)
		if len(ops) == 4 {
			if ops[3].Kind() != arch.ArmOpShift {
				return nil, fmt.Errorf("%s: shift modifier expected", name)
			}

			a := shiftAmt(ops[3])
			if a < 0 || a > 63 {
				return nil, fmt.Errorf("%s: bad shift", name)
			}

			shift, amt = ops[3].ShiftName(), uint32(a)
		}

		return makeLogShiftStruct(name, rd, rn, rm, shift, amt, rd[0] == 'x')
	}

	return nil, fmt.Errorf("%s: bad operand", name)
}

// makeLogImmStruct/makeLogShiftStruct — assembly by the base name.
func makeLogImmStruct(name, rd, rn string, immr, imms uint32, n, is64 bool) (Instr, error) {
	switch name {
	case "orr":
		return OrrImmOf(rd, rn, immr, imms, n, is64)
	case "eor":
		return EorImmOf(rd, rn, immr, imms, n, is64)
	case "ands":
		return AndsImmOf(rd, rn, immr, imms, n, is64)
	default:
		return AndImmOf(rd, rn, immr, imms, n, is64)
	}
}

func makeLogShiftStruct(name, rd, rn, rm, shift string, amt uint32, isf bool) (Instr, error) {
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
		return AndsShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	case "bic":
		return BicShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	case "orr":
		return OrrShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	case "orn":
		return OrnShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	case "eor":
		return EorShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	case "eon":
		return EonShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	case "bics":
		return BicsShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
	default:
		return AndShiftOf(common.rd, common.rn, common.rm, common.imm6, common.shift, common.isf)
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
