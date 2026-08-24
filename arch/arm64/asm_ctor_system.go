package arm64

// Assembler constructors: adr/adrp, b.cond, tbz/tbnz, system ones
// (svc/brk/hlt/hvc/udf, mrs/msr), adr points. The imm structs accept
// concrete numbers — the constructor validates the form.

import (
	"errors"
	"fmt"
)

// newAdr — adr rd, sym|#off: the target (resolveOps: abs − addr).
// Our syntax: adr rd, #<offset from the current address> (as the formatter
// prints it); a symbolic target — legacy.
func newAdr(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].kind != armOpImm {
		return nil, errors.New("adr: want rd, #off")
	}

	rd, err := wantAReg(ops[0], "adr")
	if err != nil {
		return nil, err
	}

	v := ops[1].num
	if v < -(1<<20) || v >= 1<<20 {
		return nil, errors.New("adr: offset out of range")
	}

	return Adr{
		rd:  rd,
		off: v,
	}, nil
}

// newAdrp — adrp rd, #pageOff (the page offset as the formatter prints it;
// the address annotation is cut off by the comment).
func newAdrp(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].kind != armOpImm {
		return nil, errors.New("adrp: want rd, #pageoff")
	}

	rd, err := wantAReg(ops[0], "adrp")
	if err != nil {
		return nil, err
	}

	v := ops[1].num
	if v < -(1<<20) || v >= 1<<20 {
		return nil, errors.New("adrp: offset out of range")
	}

	return Adrp{
		rd:  rd,
		off: v,
	}, nil
}

// newTbzArm — tbz/tbnz rt, #bit, target.
func newTbzArm(isTbnz bool) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 || ops[1].kind != armOpImm {
			return nil, errors.New("tbz: want rt, #bit, target")
		}

		rt, err := wantAReg(ops[0], "tbz")
		if err != nil {
			return nil, err
		}

		bit := ops[1].num
		if bit < 0 || bit > 63 {
			return nil, errors.New("tbz: bad bit")
		}

		if ops[2].kind != armOpImm {
			return nil, errors.New("tbz: target expected")
		}

		return Tbz{
			rt:     rt,
			bit:    uint32(bit),
			target: immNum(ops[2].num),
			isTbnz: isTbnz,
		}, nil
	}
}

// newSvc — svc #imm16.
func newSvc(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "svc")
	if err != nil {
		return nil, err
	}

	return sysImm{
		name:  "svc",
		imm16: uint32(v),
		enc:   0xD4000001,
		shift: 5,
	}, nil
}

// newBrkArm — brk #imm16 (#0 → "#0").
func newBrkArm(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "brk")
	if err != nil {
		return nil, err
	}

	return sysImm{
		name:  "brk",
		imm16: uint32(v),
		enc:   0xD4200000,
		shift: 5,
	}, nil
}

// newUdfArm — udf #imm16.
func newUdfArm(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "udf")
	if err != nil {
		return nil, err
	}

	return sysImm{
		name:  "udf",
		imm16: uint32(v),
		enc:   0x00000000,
		shift: 5,
	}, nil
}

// newHlt — hlt #imm16; newHvc — hvc.
func newHlt(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "hlt")
	if err != nil {
		return nil, err
	}

	return sysImm{
		name:  "hlt",
		imm16: uint32(v),
		enc:   0xD4400000,
		shift: 21,
	}, nil
}

func newHvc(ops []vOp) (Instr, error) {
	v, err := sysImm16(ops, "hvc")
	if err != nil {
		return nil, err
	}

	return sysImm{
		name:  "hvc",
		imm16: uint32(v),
		enc:   0xD4000002,
		shift: 21,
	}, nil
}

// sysImm16 — a single immediate operand.
func sysImm16(ops []vOp, name string) (int64, error) {
	if len(ops) != 1 || ops[0].kind != armOpImm {
		return 0, fmt.Errorf("%s: want #imm", name)
	}

	v := ops[0].num
	if v < 0 || v > 0xffff {
		return 0, fmt.Errorf("%s: bad imm", name)
	}

	return v, nil
}

// newMrsArm — mrs rd, sysreg.
func newMrsArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].sym == "" {
		return nil, errors.New("mrs: want rd, sysreg")
	}

	rd, err := wantAReg(ops[0], "mrs")
	if err != nil {
		return nil, err
	}

	return Mrs{
		rd:     rd,
		sysreg: ops[1].sym,
	}, nil
}

// newMsrArm — msr sysreg, rt.
func newMsrArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[0].sym == "" {
		return nil, errors.New("msr: want sysreg, rt")
	}

	rt, err := wantAReg(ops[1], "msr")
	if err != nil {
		return nil, err
	}

	return Msr{
		rt:     rt,
		sysreg: ops[0].sym,
	}, nil
}

// newBcondOf — b.<cond> target: a numeric target — the absolute; a
// symbolic one — legacy.
func newBcondOf(cond string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 1 || ops[0].kind != armOpImm {
			return nil, fmt.Errorf("b.%s: want target", cond)
		}

		v := ops[0].num

		return Bcond{
			cond:   cond,
			target: immNum(v),
		}, nil
	}
}
