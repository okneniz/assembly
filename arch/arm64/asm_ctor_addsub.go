package arm64

// Add/sub-family assembler constructors: add/adds/sub/subs (imm12
// [lsl #12], register shifts, extensions). Symbolic immediates return an
// error — they fall through to the legacy candidates; the encoding choice
// is confirmed by encodeARM's self-verify. The aliases cmp/cmn (Rd = zr),
// neg/negs (Rn = zr) — asm/arm64/alias.

import (
	"fmt"
)

// asAddSub — the constructor by the base (op, S) pair and form.
func asAddSub(ops []vOp, name string) (Instr, error) {
	if len(ops) != 3 && len(ops) != 4 {
		return nil, fmt.Errorf("%s: want rd, rn, op", name)
	}

	rd, rn, err := armReg2(ops[:2], name)
	if err != nil {
		return nil, err
	}

	rdN, rnN, err := regNums2(rd, rn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	return addSubThird(ops, name, rdN, rnN, 2)
}

// addSubThird — the third operand (imm | reg | reg+modifier) selects the form.
// idx — the index of the third operand in ops (2 for add, 1 for cmp).
func addSubThird(ops []vOp, base string, rdN, rnN uint32, idx int) (Instr, error) {
	if idx >= len(ops) {
		return nil, fmt.Errorf("%s: missing operand", base)
	}

	op := ops[idx]
	mod := vOp{}
	hasMod := idx+1 < len(ops)
	if hasMod {
		mod = ops[idx+1]
	}

	isf := ops[0].reg[0] == 'x'
	switch op.kind {
	case armOpImm:
		if hasMod {
			if mod.kind != armOpShift || mod.shift != "lsl" {
				return nil, fmt.Errorf("%s: only lsl #12 modifier allowed", base)
			}

			if shiftAmt(mod) != 12 {
				return nil, fmt.Errorf("%s: shift must be lsl #12", base)
			}
		}

		if op.sym != "" {
			return nil, fmt.Errorf("%s: want immediate operand", base)
		}

		v := op.num
		if v < 0 || v > 0xfff {
			return nil, fmt.Errorf("%s: large imm — legacy path", base)
		}

		sh := hasMod
		return makeImmStruct(base, rdN, rnN, uint32(v), sh, isf), nil
	case armOpReg:
		rm, err := wantAReg(op, base)
		if err != nil {
			return nil, err
		}

		if !hasMod {
			return makeShiftStruct(base, rdN, rnN, rm, "lsl", 0, isf), nil
		}

		switch mod.kind {
		case armOpShift:
			amt := shiftAmt(mod)
			if amt < 0 || amt > 63 {
				return nil, fmt.Errorf("%s: bad shift", base)
			}

			return makeShiftStruct(base, rdN, rnN, rm, mod.shift, uint32(amt), isf), nil
		case armOpExtend:
			ext := mod.shift
			amt := mod.num
			if amt < 0 || amt > 7 {
				return nil, fmt.Errorf("%s: bad extend amount", base)
			}

			return makeExtStruct(base, rdN, rnN, rm, ext, uint32(amt), isf), nil
		default:
			return nil, fmt.Errorf("%s: bad modifier", base)
		}
	default:
		return nil, fmt.Errorf("%s: bad operand", base)
	}
}

// makeImmStruct/makeShiftStruct/makeExtStruct — struct assembly by the base name.
func makeImmStruct(base string, rd, rn, imm12 uint32, sh, isf bool) Instr {
	switch base {
	case "adds":
		return AddsImm{
			rdNum: rd,
			rnNum: rn,
			imm12: imm12,
			shift: sh,
			isf:   isf,
		}
	case "sub":
		return SubImm{
			rdNum: rd,
			rnNum: rn,
			imm12: imm12,
			shift: sh,
			isf:   isf,
		}
	case "subs":
		return SubsImm{
			rdNum: rd,
			rnNum: rn,
			imm12: imm12,
			shift: sh,
			isf:   isf,
		}
	default:
		return AddImm{
			rdNum: rd,
			rnNum: rn,
			imm12: imm12,
			shift: sh,
			isf:   isf,
		}
	}
}

func makeShiftStruct(base string, rdN, rnN uint32, rm, shift string, amt uint32, isf bool) Instr {
	rd := addSubRegName(rdN, isf, base == "adds" || base == "subs")
	rn := addSubRegName(rnN, isf, false)
	switch base {
	case "adds":
		return AddsShift{
			rd:    rd,
			rn:    rn,
			rm:    rm,
			imm6:  amt,
			shift: shift,
			isf:   isf,
		}
	case "sub":
		return SubShift{
			rd:    rd,
			rn:    rn,
			rm:    rm,
			imm6:  amt,
			shift: shift,
			isf:   isf,
		}
	case "subs":
		return SubsShift{
			rd:    rd,
			rn:    rn,
			rm:    rm,
			imm6:  amt,
			shift: shift,
			isf:   isf,
		}
	default:
		return AddShift{
			rd:    rd,
			rn:    rn,
			rm:    rm,
			imm6:  amt,
			shift: shift,
			isf:   isf,
		}
	}
}

func makeExtStruct(base string, rdN, rnN uint32, rm, ext string, amt uint32, isf bool) Instr {
	rmN, err := armRegNum(rm)
	if err != nil {
		rmN = 0
	}

	eb := newExtBase(rdN, rnN, rmN, ext, amt, isf)
	switch base {
	case "adds":
		return AddsExt{extBase: eb}
	case "sub":
		return SubExt{extBase: eb}
	case "subs":
		return SubsExt{extBase: eb}
	default:
		return AddExt{extBase: eb}
	}
}

func newAddArm(ops []vOp) (Instr, error) {
	return asAddSub(ops, "add")
}
func newAddsArm(ops []vOp) (Instr, error) {
	return asAddSub(ops, "adds")
}
func newSubArm(ops []vOp) (Instr, error) {
	return asAddSub(ops, "sub")
}
func newSubsArm(ops []vOp) (Instr, error) {
	return asAddSub(ops, "subs")
}
