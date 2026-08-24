package arm64

// Load/store pair assembler constructors: ldp/stp rt, rt2, [rn{, #imm7}{!}]
// | ldp rt, rt2, [rn], #imm.

import (
	"fmt"
)

// newLdpArm — ldp rt, rt2, [rn{, #imm7}{!}] | ldp rt, rt2, [rn], #imm.
func newLdpArm(ops []vOp) (Instr, error) {
	return makePairCtor(ops, "ldp")
}

func newStpArm(ops []vOp) (Instr, error) {
	return makePairCtor(ops, "stp")
}

func makePairCtor(ops []vOp, name string) (Instr, error) {
	if len(ops) != 3 {
		return nil, fmt.Errorf("%s: want rt, rt2, mem", name)
	}

	rt, rt2, err := armReg2(ops, name)
	if err != nil {
		return nil, err
	}

	m := ops[2].mem
	if m == nil {
		return nil, fmt.Errorf("%s: memory operand expected", name)
	}

	rn := m.base
	scale := uint32(2)
	enc := uint32(0xA9400000) // ldp x pair offset form
	switch rt[0] {
	case 'd':
		scale, enc = 3, 0x6C400000
	case 's':
		enc = 0x2C400000
	case 'w':
		scale, enc = 2, 0x28400000
	}

	if name == "stp" {
		enc &^= 1 << 22
	}

	var kind memKind
	var off int64
	switch {
	case m.post != 0:
		kind, off = memPost, m.post
		enc = enc&^0x01800000 | 0x00800000
	case m.pre:
		kind, off = memPre, m.off
		enc |= 0x01800000
	default:
		kind, off = memImm, m.off
	}

	pb := newPairBase(rt, rt2, rn, kind, off, scale, enc)
	if name == "ldp" {
		return Ldp{pairBase: pb}, nil
	}

	return Stp{pairBase: pb}, nil
}
