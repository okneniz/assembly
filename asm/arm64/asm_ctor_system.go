package arm64

// Assembler constructors: adr/adrp, b.cond, tbz/tbnz, system ones
// (svc/brk/hlt/hvc/udf, mrs/msr), adr points. The imm structs accept
// concrete numbers — the constructor validates the form.

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// newAdr — adr rd, sym|#off: the target (resolveOps: abs − addr).
// Our syntax: adr rd, #<offset from the current address> (as the formatter
// prints it); a symbolic target — legacy.
func newAdr(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].Kind() != arch.ArmOpImm {
		return nil, errors.New("adr: want rd, #off")
	}

	rd, err := wantAReg(ops[0], "adr")
	if err != nil {
		return nil, err
	}

	v := ops[1].Num()
	if v < -(1<<20) || v >= 1<<20 {
		return nil, errors.New("adr: offset out of range")
	}

	return AdrOf(rd, v), nil
}

// newAdrp — adrp rd, #pageOff (the page offset as the formatter prints it;
// the address annotation is cut off by the comment).
func newAdrp(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].Kind() != arch.ArmOpImm {
		return nil, errors.New("adrp: want rd, #pageoff")
	}

	rd, err := wantAReg(ops[0], "adrp")
	if err != nil {
		return nil, err
	}

	v := ops[1].Num()
	if v < -(1<<20) || v >= 1<<20 {
		return nil, errors.New("adrp: offset out of range")
	}

	return AdrpOf(rd, v), nil
}

// newTbzArm — tbz/tbnz rt, #bit, target.
func newTbzArm(isTbnz bool) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 || ops[1].Kind() != arch.ArmOpImm {
			return nil, errors.New("tbz: want rt, #bit, target")
		}

		rt, err := wantAReg(ops[0], "tbz")
		if err != nil {
			return nil, err
		}

		bit := ops[1].Num()
		if bit < 0 || bit > 63 {
			return nil, errors.New("tbz: bad bit")
		}

		if ops[2].Kind() != arch.ArmOpImm {
			return nil, errors.New("tbz: target expected")
		}

		return TbzOf(rt, uint32(bit), ops[2].Num(), isTbnz), nil
	}
}

// newSvc — svc #imm16.
func newSvc(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "svc")
	if err != nil {
		return nil, err
	}

	return arch.SysImmOf("svc", uint32(v), 0xD4000001, 5), nil
}

// newBrkArm — brk #imm16 (#0 → "#0").
func newBrkArm(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "brk")
	if err != nil {
		return nil, err
	}

	return arch.SysImmOf("brk", uint32(v), 0xD4200000, 5), nil
}

// newUdfArm — udf #imm16.
func newUdfArm(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "udf")
	if err != nil {
		return nil, err
	}

	return arch.SysImmOf("udf", uint32(v), 0x00000000, 5), nil
}

// newHlt — hlt #imm16; newHvc — hvc.
func newHlt(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "hlt")
	if err != nil {
		return nil, err
	}

	return arch.SysImmOf("hlt", uint32(v), 0xD4400000, 21), nil
}

func newHvc(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "hvc")
	if err != nil {
		return nil, err
	}

	return arch.SysImmOf("hvc", uint32(v), 0xD4000002, 21), nil
}

// sysImm16 — a single immediate operand.
func sysImm16(ops []vOp, name string) (int64, error) {
	if len(ops) != 1 || ops[0].Kind() != arch.ArmOpImm {
		return 0, fmt.Errorf("%s: want #imm", name)
	}

	v := ops[0].Num()
	if v < 0 || v > 0xffff {
		return 0, fmt.Errorf("%s: bad imm", name)
	}

	return v, nil
}

// newMrsArm — mrs rd, sysreg.
func newMrsArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].Sym() == "" {
		return nil, errors.New("mrs: want rd, sysreg")
	}

	rd, err := wantAReg(ops[0], "mrs")
	if err != nil {
		return nil, err
	}

	return MrsOf(rd, ops[1].Sym()), nil
}

// newMsrArm — msr sysreg, rt.
func newMsrArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[0].Sym() == "" {
		return nil, errors.New("msr: want sysreg, rt")
	}

	rt, err := wantAReg(ops[1], "msr")
	if err != nil {
		return nil, err
	}

	return MsrOf(rt, ops[0].Sym()), nil
}

// newBcondOf — b.<cond> target: a numeric target — the absolute; a
// symbolic one — legacy.
func newBcondOf(cond string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 1 || ops[0].Kind() != arch.ArmOpImm {
			return nil, fmt.Errorf("b.%s: want target", cond)
		}

		v := ops[0].Num()

		return BcondOf(cond, v), nil
	}
}
