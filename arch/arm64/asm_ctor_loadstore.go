package arm64

// Load/store assembler constructors: ldr/str (imm12 scale, unscaled
// ldur/stur, post/pre-index, register offset), signed ldrs*, pairs
// ldp/stp, atomics ldar/stlr/ldaxr, exclusives stlxr/stxrb. Literal/label
// arrive computed (resolveOps of the syntax layer).

import (
	"errors"
	"fmt"
)

// memOperand — parse a mem operand: base, addressing kind, offset/index.
func memOperand(
	op vOp,
	name string,
) (base string, kind memKind, off int64, rm, option string, shiftAmt32 uint32, err error) {
	if op.mem == nil {
		return "", 0, 0, "", "", 0, fmt.Errorf("%s: memory operand expected", name)
	}

	m := op.mem
	base = m.base
	if m.post != 0 {
		return base, memPost, m.post, "", "", 0, nil
	}

	if m.offReg != "" {
		if m.optAmt < 0 || m.optAmt > 7 {
			return "", 0, 0, "", "", 0, fmt.Errorf("%s: bad option amount", name)
		}

		opt := m.opt
		if opt == "" {
			opt = "lsl"
		}

		return base, memRegOff, 0, m.offReg, opt, uint32(m.optAmt), nil
	}

	if m.pre {
		return base, memPre, m.off, "", "", 0, nil
	}

	return base, memUnscaled, m.off, "", "", 0, nil // imm12 vs unscaled — chosen by divisibility
}

// Constructors by mnemonic (enc from the decode table).
func newLdrArm(ops []vOp) (Instr, error) {
	// literal pool: ldr rt, =literal
	if len(ops) == 2 && ops[1].kind == armOpLit {
		return newLdrPool(ops[0], ops[1])
	}

	// literal form: ldr rt, label (no brackets)
	if len(ops) == 2 && ops[1].mem == nil {
		return newLdrLiteral(ops)
	}

	return makeLSCtorTyped(ops, "ldr", 0xF9400000)
}
func newLdrbArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldrb", 0x39400000)
}
func newLdrhArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldrh", 0x79400000)
}
func newStrArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "str", 0xF9000000)
}
func newStrbArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "strb", 0x39000000)
}
func newStrhArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "strh", 0x79000000)
}
func newLdursArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldur", 0xF8400000)
}
func newStursArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "stur", 0xF8000000)
}
func newLdurbArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldurb", 0x38400000)
}
func newLdurhArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldurh", 0x78400000)
}
func newSturbArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "sturb", 0x38000000)
}
func newSturhArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "sturh", 0x78000000)
}
func newLdrswArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldrsw", 0xB9800000)
}
func newLdrsbArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldrsb", 0x39800000)
}
func newLdrshArm(ops []vOp) (Instr, error) {
	return makeLSCtorTyped(ops, "ldrsh", 0x79800000)
}

// makeLSBase — assemble an lsBase.
func makeLSBase(
	rt, rn string,
	kind memKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) lsBase {
	return lsBase{
		rt:       rt,
		rn:       rn,
		kind:     kind,
		off:      off,
		enc:      enc,
		rm:       rm,
		option:   option,
		shiftAmt: amt,
	}
}

// makeLSCtorTyped — builds the struct by mnemonic.
func makeLSCtorTyped(ops []vOp, name string, enc uint32) (Instr, error) {
	if len(ops) != 2 {
		return nil, fmt.Errorf("%s: want rt, [rn{, off}]", name)
	}

	rt, err := wantAReg(ops[0], name)
	if err != nil {
		if ops[0].reg == "" {
			return nil, fmt.Errorf("%s: register expected", name)
		}

		rt = ops[0].reg
	}

	rn, kind, off, rm, option, amt, err := memOperand(ops[1], name)
	if err != nil {
		return nil, err
	}

	scale := enc >> 30 & 3
	if kind == memUnscaled && off >= 0 && off&(int64(1)<<scale-1) == 0 {
		kind = memImm
	}

	switch name {
	case "ldr":
		return Ldr{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "ldrb":
		return Ldrb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "ldrh":
		return Ldrh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "str":
		return Str{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "strb":
		return Strb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "strh":
		return Strh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "ldur":
		return Ldur{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "stur":
		return Stur{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "ldurb":
		return Ldurb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "ldurh":
		return Ldurh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "sturb":
		return Sturb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "sturh":
		return Sturh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	case "ldrsw":
		return Ldrsw{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}, nil
	}

	// ldrsb/ldrsh — own fields (imm12 scale), not lsBase
	if name == "ldrsb" {
		return Ldrsb{
			rt:  rt,
			rn:  rn,
			off: off,
		}, nil
	}

	return Ldrsh{
		rt:  rt,
		rn:  rn,
		off: off,
	}, nil
}

// newLdrPool — ldr rt, =literal: a literal-pool slot at the end of the
// subsection (see asm.PoolUser). Unlike GAS (which emits movz/movk for
// simple literals) — always ldr-literal + pool: the semantics are
// equivalent, the bytes differ. The target — the slot's auto-name
// (PoolName), resolved by the core.
func newLdrPool(rt vOp, lit vOp) (Instr, error) {
	if rt.reg == "" {
		return nil, errors.New("ldr: register expected")
	}

	if _, err := armRegNum(rt.reg); err != nil {
		return nil, fmt.Errorf("ldr: %w", err)
	}

	enc := uint32(0x58000000)
	if rt.reg[0] == 'w' {
		enc = 0x18000000
	}

	lb := newLsBase(rt.reg, "", memLiteral, 0, uint64(lit.num), enc, "", "", 0)
	return ldrPoolWrap{Ldr{lsBase: lb}}, nil
}

// newLdrLiteral — ldr rt, label|#addr: the target is already computed
// (resolveOps); w/x by the rt type.
func newLdrLiteral(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].kind != armOpImm {
		return nil, errors.New("ldr: want rt, target")
	}

	rt := ops[0].reg
	if rt == "" {
		return nil, errors.New("ldr: register expected")
	}

	if _, err := armRegNum(rt); err != nil {
		return nil, fmt.Errorf("ldr: %w", err)
	}

	// numeric target — the absolute (in tgt)
	enc := uint32(0x58000000) // x form
	if rt[0] == 'w' {
		enc = 0x18000000
	}

	lb := newLsBase(rt, "", memLiteral, 0, uint64(ops[1].num), enc, "", "", 0)
	return Ldr{lsBase: lb}, nil
}
