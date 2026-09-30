package arm64

import (
	"errors"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// Assembler constructors for the IC/DC/TLBI system operations and the
// eret family: the keyword operand survives resolution through the arch
// keyword table, the optional Xt register joins the word.

// newSysOpArm - ic/dc/tlbi <op>{, xt}: the first operand is the
// operation spelling (a keyword symbol), the second - when the
// operation takes one - the Xt register.
func newSysOpArm(mnem string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) == 0 || ops[0].Sym() == "" {
			return nil, errors.New(mnem + ": want operation")
		}

		rt := ""
		if len(ops) == 2 {
			if !ops[1].IsReg() {
				return nil, errors.New(mnem + ": want op, xt")
			}

			rt = ops[1].Reg()
		}

		if len(ops) > 2 {
			return nil, errors.New(mnem + ": want op[, xt]")
		}

		return arch.SysOpOf(mnem, ops[0].Sym(), rt)
	}
}

// newEret - eret/eretaa/eretab: a fixed word, no operands.
func newEret(name string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 0 {
			return nil, errors.New(name + " expects no operands")
		}

		return arch.SysFixedOf(name, "", "Pseudo", enc), nil
	}
}
