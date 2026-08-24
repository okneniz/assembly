package alias

// Aliases of the add/sub family: cmp/cmn (base subs/adds with Rd = 31,
// zr when S=1), neg/negs (base sub/subs with Rn = zr). The form is
// chosen by the arch family machinery (AddSubThird); symbolic immediates
// return an error and fall through to the legacy candidates, as with the
// base instructions.

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// newCmp is the cmp/cmn rn, op alias: base subs/adds with Rd = 31 (zr
// when S=1).
func newCmp(base string) arch.ArmCtor {
	return func(ops []arch.ArmOp) (arch.Instr, error) {
		if len(ops) != 2 && len(ops) != 3 {
			return nil, fmt.Errorf("%s: want rn, op", base)
		}

		rn, err := arch.WantAReg(ops[0], base)
		if err != nil {
			return nil, err
		}

		rnN, err := arch.ArmRegNum(rn)
		if err != nil {
			return nil, err
		}

		return arch.AddSubThird(ops, base, 31, rnN, 1)
	}
}

// newNeg is the neg/negs rd, rm{, shift #imm6} alias: base sub/subs
// with Rn = zr. The third operand is an optional shift; without it
// lsl #0 (canonical).
func newNeg(base string) arch.ArmCtor {
	return func(ops []arch.ArmOp) (arch.Instr, error) {
		if len(ops) != 2 && len(ops) != 3 {
			return nil, fmt.Errorf("%s: want rd, rm", base)
		}

		rd, rm, err := arch.ArmReg2(ops[:2], base)
		if err != nil {
			return nil, err
		}

		isf := rd[0] == 'x'
		sh, imm6 := "lsl", uint32(0)
		if len(ops) == 3 {
			op := ops[2]
			if !op.IsShift() {
				return nil, fmt.Errorf("%s: want shift, got operand %d", base, 3)
			}

			sh = op.ShiftName()
			imm6 = uint32(arch.ShiftAmt(op))
		}

		if base == "sub" {
			return arch.SubShiftOf(rd, arch.ZeroReg(rd), rm, imm6, sh, isf), nil
		}

		return arch.SubsShiftOf(rd, arch.ZeroReg(rd), rm, imm6, sh, isf), nil
	}
}
