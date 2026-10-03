package arm64

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// Load/store pair assembler constructors: ldp/stp rt, rt2, [rn{, #imm7}{!}]
// | ldp rt, rt2, [rn], #imm.

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

	if !ops[2].IsMem() {
		return nil, fmt.Errorf("%s: memory operand expected", name)
	}

	m := ops[2].Mem()

	rn := m.Base()
	scale := uint32(2)
	enc := uint32(0x29400000) // ldp w pair offset form
	switch rt[0] {
	case 'x':
		scale, enc = 3, 0xA9400000
	case 's':
		enc = 0x2D400000
	case 'd':
		scale, enc = 3, 0x6D400000
	case 'q':
		scale, enc = 4, 0xAD400000
	case 'w': // the w defaults stand
	default:
		return nil, fmt.Errorf("%s: x/w/s/d/q register expected, got %q", name, rt)
	}

	if rt2[0] != rt[0] {
		return nil, fmt.Errorf(
			"%s: pair registers must share the width, got %q and %q",
			name,
			rt,
			rt2,
		)
	}

	if name == "stp" {
		enc &^= 1 << 22
	}

	var kind arch.MemKind
	var off int64
	switch {
	case m.Post() != 0:
		kind, off = arch.MemPost, m.Post()
		enc = enc&^0x01800000 | 0x00800000
	case m.Pre():
		kind, off = arch.MemPre, m.Off()
		enc |= 0x01800000
	default:
		kind, off = arch.MemImm, m.Off()
	}

	if off&(int64(1)<<scale-1) != 0 || off>>scale < -64 || off>>scale > 63 {
		return nil, fmt.Errorf(
			"%s: offset %d is not a multiple of %d or out of the imm7 range",
			name,
			off,
			1<<scale,
		)
	}

	if name == "ldp" {
		return LdpOf(rt, rt2, rn, kind, off, scale, enc)
	}

	return StpOf(rt, rt2, rn, kind, off, scale, enc)
}
