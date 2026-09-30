package arm64

import (
	"errors"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// Assembler constructors for hints and system hints: the barrier
// family, yield, prfm pldl1keep.

// newBarrierArm — dmb/dsb option, isb [sy]: the text spelling becomes
// the typed domain (dmb/dsb choose one, isb has none to choose).
func newBarrierArm(name string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		var b arch.Builder
		switch len(ops) {
		case 0:
			if name != "isb" {
				return nil, errors.New(name + ": want option")
			}

			return b.Isb()
		case 1:
			domain, ok := arch.DomainOf(ops[0].Sym())
			if !ok {
				return nil, errors.New(name + ": unknown option " + ops[0].Sym())
			}

			if name == "isb" {
				if domain != arch.Sy {
					return nil, errors.New(name + ": takes sy only")
				}

				return b.Isb()
			}

			if name == "dmb" {
				return b.Dmb(domain)
			}

			return b.Dsb(domain)
		default:
			return nil, errors.New(name + ": want option")
		}
	}
}

// newHint — the operandless aliases of the HINT space (yield/wfi/wfe/
// sev/sevl); the word is fixed.
func newHint(name string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 0 {
			return nil, errors.New(name + " expects no operands")
		}

		return arch.SysFixedOf(name, "", "Hint", enc), nil
	}
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

	return PrfmOf(ops[1].Mem().Base())
}
