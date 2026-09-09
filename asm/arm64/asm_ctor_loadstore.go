package arm64

// Load/store assembler constructors: ldr/str (imm12 scale, unscaled
// ldur/stur, post/pre-index, register offset), signed ldrs*, pairs
// ldp/stp, atomics ldar/stlr/ldaxr, exclusives stlxr/stxrb. Literal/label
// arrive computed (resolveOps of the syntax layer).

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// memOperand — parse a mem operand: base, addressing kind, offset/index.
func memOperand(
	op vOp,
	name string,
) (base string, kind arch.MemKind, off int64, rm, option string, shiftAmt32 uint32, err error) {
	if !op.IsMem() {
		return "", 0, 0, "", "", 0, fmt.Errorf("%s: memory operand expected", name)
	}

	m := op.Mem()
	base = m.Base()
	if m.Post() != 0 {
		return base, arch.MemPost, m.Post(), "", "", 0, nil
	}

	if m.OffReg() != "" {
		if m.OptAmt() < 0 || m.OptAmt() > 7 {
			return "", 0, 0, "", "", 0, fmt.Errorf("%s: bad option amount", name)
		}

		opt := m.Opt()
		if opt == "" {
			opt = "lsl"
		}

		return base, arch.MemRegOff, 0, m.OffReg(), opt, uint32(m.OptAmt()), nil
	}

	if m.Pre() {
		return base, arch.MemPre, m.Off(), "", "", 0, nil
	}

	return base, arch.MemUnscaled, m.Off(), "", "", 0, nil // imm12 vs unscaled — chosen by divisibility
}

// Constructors by mnemonic (enc from the decode table).
func newLdrArm(ops []vOp) (Instr, error) {
	// literal pool: ldr rt, =literal
	if len(ops) == 2 && ops[1].Kind() == arch.ArmOpLit {
		return newLdrPool(ops[0], ops[1])
	}

	// literal form: ldr rt, label (no brackets)
	if len(ops) == 2 && !ops[1].IsMem() {
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

// makeLSCtorTyped — builds the struct by mnemonic.
func makeLSCtorTyped(ops []vOp, name string, enc uint32) (Instr, error) {
	if len(ops) != 2 {
		return nil, fmt.Errorf("%s: want rt, [rn{, off}]", name)
	}

	rt, err := wantAReg(ops[0], name)
	if err != nil {
		if ops[0].Reg() == "" {
			return nil, fmt.Errorf("%s: register expected", name)
		}

		rt = ops[0].Reg()
	}

	rn, kind, off, rm, option, amt, err := memOperand(ops[1], name)
	if err != nil {
		return nil, err
	}

	scale := enc >> 30 & 3
	if kind == arch.MemUnscaled && off >= 0 && off&(int64(1)<<scale-1) == 0 {
		kind = arch.MemImm
	}

	switch name {
	case "ldr":
		return LdrOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "ldrb":
		return LdrbOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "ldrh":
		return LdrhOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "str":
		return StrOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "strb":
		return StrbOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "strh":
		return StrhOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "ldur":
		return LdurOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "stur":
		return SturOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "ldurb":
		return LdurbOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "ldurh":
		return LdurhOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "sturb":
		return SturbOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "sturh":
		return SturhOf(rt, rn, kind, off, enc, rm, option, amt), nil
	case "ldrsw":
		return LdrswOf(rt, rn, kind, off, enc, rm, option, amt), nil
	}

	// ldrsb/ldrsh — own fields (imm12 scale), not lsBase
	if name == "ldrsb" {
		return LdrsbOf(rt, rn, off), nil
	}

	return LdrshOf(rt, rn, off), nil
}

// newLdrPool — ldr rt, =literal: a literal-pool slot at the end of the
// subsection (see asm.PoolUser). Unlike GAS (which emits movz/movk for
// simple literals) — always ldr-literal + pool: the semantics are
// equivalent, the bytes differ. The target — the slot's auto-name
// (PoolName), resolved by the core.
func newLdrPool(rt vOp, lit vOp) (Instr, error) {
	if rt.Reg() == "" {
		return nil, errors.New("ldr: register expected")
	}

	if _, err := armRegNum(rt.Reg()); err != nil {
		return nil, fmt.Errorf("ldr: %w", err)
	}

	enc := uint32(0x58000000)
	if rt.Reg()[0] == 'w' {
		enc = 0x18000000
	}

	return arch.LdrPoolWrapOf(rt.Reg(), lit.Num(), enc), nil
}

// newLdrLiteral — ldr rt, label|#addr: the target is already computed
// into a pc-relative offset (resolveOps); w/x by the rt type.
func newLdrLiteral(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].Kind() != arch.ArmOpImm {
		return nil, errors.New("ldr: want rt, target")
	}

	rt := ops[0].Reg()
	if rt == "" {
		return nil, errors.New("ldr: register expected")
	}

	if _, err := armRegNum(rt); err != nil {
		return nil, fmt.Errorf("ldr: %w", err)
	}

	enc := uint32(0x58000000) // x form
	if rt[0] == 'w' {
		enc = 0x18000000
	}

	return LdrLitOf(rt, ops[1].Num(), enc), nil
}
