package arm64

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

	return sysFixed{
		name:  "dmb",
		ops:   "st",
		group: "System",
		enc:   0xD5033EBF,
	}, nil
}

// newYield — yield (no operands).
func newYield(ops []vOp) (Instr, error) {
	if len(ops) != 0 {
		return nil, errors.New("yield expects no operands")
	}

	return sysFixed{
		name:  "yield",
		ops:   "",
		group: "Hint",
		enc:   0xD503203F,
	}, nil
}

// newDc — dc zva, x0 (fixed).
func newDc(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("dc: want zva, x0")
	}

	return sysFixed{
		name:  "dc",
		ops:   "zva, x0",
		group: "System",
		enc:   0xD50B7420,
	}, nil
}

// newPrfmArm — prfm pldl1keep, [rn] (fixed form).
func newPrfmArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[1].mem == nil {
		return nil, errors.New("prfm: want op, [rn]")
	}

	// the first operand — the pldl1keep keyword (resolveOps kept the name)
	if ops[0].sym != "pldl1keep" {
		return nil, errors.New("prfm: pldl1keep expected")
	}

	return Prfm{rn: ops[1].mem.base}, nil
}
