package alias

// R-form aliases: mul/mneg (base madd/msub with ra = zr) and the
// inverted aliases of the csel family - cset/csetm/cinc/cinv/cneg
// carry an inverted condition (that is how ObjDump prints them from
// the base csinc/csinv/csneg encoding).

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// newMul is the mul rd, rn, rm alias: madd with ra = zr.
func newMul(ops []arch.ArmOp) (arch.Instr, error) {
	return arch.MaddCtor(ops, "mul")
}

// newMneg is the mneg rd, rn, rm alias: msub with ra = zr.
func newMneg(ops []arch.ArmOp) (arch.Instr, error) {
	return arch.Msub3(ops, "mneg")
}

// newCset/newCsetm are the cset/csetm rd, cond aliases: csinc/csinv
// with Rd=Rn=Rm=zr.
func newCset(ops []arch.ArmOp) (arch.Instr, error) {
	return csCond(ops, "cset", "csinc")
}

func newCsetm(ops []arch.ArmOp) (arch.Instr, error) {
	return csCond(ops, "csetm", "csinv")
}

func csCond(ops []arch.ArmOp, name, base string) (arch.Instr, error) {
	if len(ops) != 2 { // cset/csetm rd, cond
		return nil, fmt.Errorf("%s: bad operands", name)
	}

	rd, err := arch.WantAReg(ops[0], name)
	if err != nil {
		return nil, err
	}

	cond, err := arch.WantCond(ops[1], name)
	if err != nil {
		return nil, err
	}

	zr := arch.ZeroReg(rd)
	return arch.CsOf(base, rd, zr, zr, arch.InvertCond(cond))
}

// newCinc/newCinv/newCneg are the cinc/cinv/cneg rd, rm, cond aliases:
// csinc/csinv/csneg with Rn=Rm.
func newCinc(ops []arch.ArmOp) (arch.Instr, error) {
	return csRM(ops, "cinc", "csinc")
}

func newCinv(ops []arch.ArmOp) (arch.Instr, error) {
	return csRM(ops, "cinv", "csinv")
}

func newCneg(ops []arch.ArmOp) (arch.Instr, error) {
	return csRM(ops, "cneg", "csneg")
}

func csRM(ops []arch.ArmOp, name, base string) (arch.Instr, error) {
	if len(ops) != 3 { // cinc/cinv/cneg rd, rm, cond
		return nil, fmt.Errorf("%s: bad operands", name)
	}

	rd, err := arch.WantAReg(ops[0], name)
	if err != nil {
		return nil, err
	}

	rm, err := arch.WantAReg(ops[1], name)
	if err != nil {
		return nil, err
	}

	cond, err := arch.WantCond(ops[2], name)
	if err != nil {
		return nil, err
	}

	return arch.CsOf(base, rd, rm, rm, arch.InvertCond(cond))
}
