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

	if name == "ldp" {
		return LdpOf(rt, rt2, rn, kind, off, scale, enc)
	}

	return StpOf(rt, rt2, rn, kind, off, scale, enc)
}
