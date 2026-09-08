package arm64

import arch "github.com/okneniz/assembly/arch/arm64"

// Assembler constructors for hints and system hints: dmb st, yield,
// dc zva, prfm pldl1keep (fixed forms, like the decode table).

import (
	"errors"
)

// newDmb — dmb st (fixed form, like the decode table).
func newDmb(ops []vOp) (Instr, error) {
	if len(ops) != 1 {
		return nil, errors.New("dmb: want option")
	}

	return arch.SysFixedOf("dmb", "st", "System", 0xD5033EBF), nil
}

// newYield — yield (no operands).
func newYield(ops []vOp) (Instr, error) {
	if len(ops) != 0 {
		return nil, errors.New("yield expects no operands")
	}

	return arch.SysFixedOf("yield", "", "Hint", 0xD503203F), nil
}

// newDc — dc zva, x0 (fixed).
func newDc(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("dc: want zva, x0")
	}

	return arch.SysFixedOf(
		"dc", "zva, x0",
		"System", 0xD50B7420), nil
}

// newPrfmArm — prfm pldl1keep, [rn] (fixed form).
func newPrfmArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 || !ops[1].IsMem() {
		return nil, errors.New("prfm: want op, [rn]")
	}

	// the first operand — the pldl1keep keyword (resolveOps kept the name)
	if ops[0].Sym() != "pldl1keep" {
		return nil, errors.New("prfm: pldl1keep expected")
	}

	return PrfmOf(ops[1].Mem().Base()), nil
}
