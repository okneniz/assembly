package arm64

import (
	"errors"
	"strings"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// Assembler constructors for the MSR (immediate) PSTATE spellings: the
// <pstatefield> operand of the real MSR_imm form (msr daifset, #1) - the
// field vocabulary is the arch keyword table (arch.PstateBase), the
// immediate rides CRm [11:8].

// pstateWord - the MSR (immediate) alias: the pstatefield operand has no
// decode-side print, so the word skips the self-verify (there are no
// free bits beyond the imm).
type pstateWord struct {
	arch.Instr
}

func (pstateWord) SkipVerify() {}

// newMsrPstate - msr <pstatefield>, #imm.
func newMsrPstate(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[0].Sym() == "" || ops[1].Kind() != arch.ArmOpImm {
		return nil, errors.New("msr: want pstate, #imm")
	}

	base, imm1, ok := arch.PstateBase(ops[0].Sym())
	if !ok {
		return nil, errors.New("msr: unknown pstate field " + ops[0].Sym())
	}

	maxImm := int64(15)
	if imm1 {
		maxImm = 1
	}

	v := ops[1].Num()
	if v < 0 || v > maxImm {
		return nil, errors.New("msr " + strings.ToLower(ops[0].Sym()) + ": bad imm")
	}

	st := arch.SysImmOf("msr "+strings.ToLower(ops[0].Sym()), uint32(v), base, 8)
	return pstateWord{Instr: st}, nil
}
