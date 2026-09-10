package arm64

// Integer R-form assembler constructors: rbit/rev/clz, register shifts
// lsl/lsr/asr/ror, udiv/sdiv/smulh/umulh, adc, ccmp, extr, madd/msub, the
// csel family. The aliases mul/mneg (ra = zr) and cset/csetm/cinc/cinv/
// cneg (inverted condition) — asm/arm64/alias. 32/64-bit width — by the
// first register (sfMatch).

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

func newRbit(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rbit: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rbit")
	if err != nil {
		return nil, err
	}

	return RbitOf(rd, rn)
}

func newRev16(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rev16: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rev16")
	if err != nil {
		return nil, err
	}

	return Rev16Of(rd, rn)
}

func newRev32(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rev32: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rev32")
	if err != nil {
		return nil, err
	}

	return Rev32Of(rd, rn)
}

func newRev(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rev: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rev")
	if err != nil {
		return nil, err
	}

	return RevOf(rd, rn)
}

func newClz(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("clz: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "clz")
	if err != nil {
		return nil, err
	}

	return ClzOf(rd, rn)
}

func newCls(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("cls: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "cls")
	if err != nil {
		return nil, err
	}

	return ClsOf(rd, rn)
}

func newLslReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("lsl: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "lsl")
	if err != nil {
		return nil, err
	}

	return LslRegOf(rd, rn, rm)
}

func newLsrReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("lsr: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "lsr")
	if err != nil {
		return nil, err
	}

	return LsrRegOf(rd, rn, rm)
}

func newAsrReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("asr: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "asr")
	if err != nil {
		return nil, err
	}

	return AsrRegOf(rd, rn, rm)
}

func newRorReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("ror: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "ror")
	if err != nil {
		return nil, err
	}

	return RorRegOf(rd, rn, rm)
}

func newUdiv(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("udiv: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "udiv")
	if err != nil {
		return nil, err
	}

	return UdivOf(rd, rn, rm)
}

func newSdiv(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("sdiv: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "sdiv")
	if err != nil {
		return nil, err
	}

	return SdivOf(rd, rn, rm)
}

func newSmulh(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("smulh: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "smulh")
	if err != nil {
		return nil, err
	}

	return SmulhOf(rd, rn, rm)
}

func newUmulh(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("umulh: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "umulh")
	if err != nil {
		return nil, err
	}

	return UmulhOf(rd, rn, rm)
}

func newAdc(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("adc: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "adc")
	if err != nil {
		return nil, err
	}

	return AdcOf(rd, rn, rm)
}

func newCcmp(ops []vOp) (Instr, error) {
	if len(ops) != 4 || ops[2].Kind() != arch.ArmOpImm || ops[3].Sym() == "" {
		return nil, errors.New("ccmp: want rn, rm, #imm, cond")
	}

	rn, rm, err := armReg2(ops, "ccmp")
	if err != nil {
		return nil, err
	}

	v := ops[2].Num()
	if v < 0 || v > 15 {
		return nil, errors.New("ccmp: bad imm")
	}

	return CcmpOf(rn, rm, uint32(v), ops[3].Sym())
}

func newExtr(ops []vOp) (Instr, error) {
	if len(ops) != 4 || ops[3].Kind() != arch.ArmOpImm {
		return nil, errors.New("extr: want rd, rn, rm, #lsb")
	}

	rd, rn, rm, err := armReg3(ops, "extr")
	if err != nil {
		return nil, err
	}

	v := ops[3].Num()
	if v < 0 || v > 63 {
		return nil, errors.New("extr: bad lsb")
	}

	return ExtrOf(rd, rn, rm, uint32(v))
}

// newMadd — madd rd, rn, rm, ra (mul/mneg with ra = zr — asm/arm64/alias).
func newMadd(ops []vOp) (Instr, error) {
	return makeMaddCtor(ops, "madd")
}

func newMsub(ops []vOp) (Instr, error) {
	if len(ops) == 4 {
		rd, rn, rm, err := armReg3(ops, "msub")
		if err != nil {
			return nil, err
		}

		ra, err := wantAReg(ops[3], "msub")
		if err != nil {
			return nil, err
		}

		return MsubOf(rd, rn, rm, ra)
	}

	return msub3(ops, "msub")
}

func msub3(ops []vOp, name string) (Instr, error) {
	if len(ops) != 3 {
		return nil, fmt.Errorf("%s: want rd, rn, rm", name)
	}

	rd, rn, rm, err := armReg3(ops, name)
	if err != nil {
		return nil, err
	}

	return MsubOf(rd, rn, rm, zeroReg(rd))
}

func makeMaddCtor(ops []vOp, name string) (Instr, error) {
	if len(ops) == 4 {
		rd, rn, rm, err := armReg3(ops, name)
		if err != nil {
			return nil, err
		}

		ra, err := wantAReg(ops[3], name)
		if err != nil {
			return nil, err
		}

		return MaddOf(rd, rn, rm, ra)
	}

	if len(ops) != 3 {
		return nil, fmt.Errorf("%s: want rd, rn, rm[, ra]", name)
	}

	rd, rn, rm, err := armReg3(ops, name)
	if err != nil {
		return nil, err
	}

	return MaddOf(rd, rn, rm, zeroReg(rd))
}

func newCselArm(ops []vOp) (Instr, error) {
	if len(ops) != 4 {
		return nil, errors.New("csel: want rd, rn, rm, cond")
	}

	rd, rn, rm, err := armReg3(ops, "csel")
	if err != nil {
		return nil, err
	}

	cond, err := wantCond(ops[3], "csel")
	if err != nil {
		return nil, err
	}

	return CselOf(rd, rn, rm, cond)
}

// csFam — the csinc/csinv/csneg constructor (full form rd, rn, rm, cond).
// The inverse aliases (cset/csetm/cinc/cinv/cneg with an inverted
// condition) — asm/arm64/alias.
func csFam(ops []vOp, name string) (Instr, error) {
	if len(ops) == 4 { // csinc/csinv/csneg rd, rn, rm, cond
		rd, rn, rm, err := armReg3(ops, name)
		if err != nil {
			return nil, err
		}

		cond, err := wantCond(ops[3], name)
		if err != nil {
			return nil, err
		}

		return csOf(name, rd, rn, rm, cond)
	}

	return nil, fmt.Errorf("%s: want rd, rn, rm, cond", name)
}

// csOf — assemble a family struct by the base encoding name.
func csOf(base string, rd, rn, rm, cond string) (Instr, error) {
	c, cerr := newCsel(rd, rn, rm, cond)
	switch base {
	case "csinc":
		return CsincOf(c)
	case "csinv":
		return CsinvOf(c)
	case "csneg":
		return CsnegOf(c)
	}

	if cerr != nil {
		return nil, cerr
	}

	return c, nil
}

func newCsinc(ops []vOp) (Instr, error) {
	return csFam(ops, "csinc")
}

func newCsinv(ops []vOp) (Instr, error) {
	return csFam(ops, "csinv")
}

func newCsneg(ops []vOp) (Instr, error) {
	return csFam(ops, "csneg")
}
