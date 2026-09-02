package arm64

// Integer R-form assembler constructors: rbit/rev/clz, register shifts
// lsl/lsr/asr/ror, udiv/sdiv/smulh/umulh, adc, ccmp, extr, madd/msub, the
// csel family. The aliases mul/mneg (ra = zr) and cset/csetm/cinc/cinv/
// cneg (inverted condition) — asm/arm64/alias. 32/64-bit width — by the
// first register (sfMatch).

import (
	"errors"
	"fmt"
)

func newRbit(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rbit: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rbit")
	if err != nil {
		return nil, err
	}

	return Rbit{
		rd: rd,
		rn: rn,
	}, nil
}

func newRev16(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rev16: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rev16")
	if err != nil {
		return nil, err
	}

	return Rev16{
		rd: rd,
		rn: rn,
	}, nil
}

func newRev32(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rev32: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rev32")
	if err != nil {
		return nil, err
	}

	return Rev32{
		rd: rd,
		rn: rn,
	}, nil
}

func newRev(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("rev: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "rev")
	if err != nil {
		return nil, err
	}

	return Rev{
		rd: rd,
		rn: rn,
	}, nil
}

func newClz(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("clz: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "clz")
	if err != nil {
		return nil, err
	}

	return Clz{
		rd: rd,
		rn: rn,
	}, nil
}

func newCls(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("cls: want rd, rn")
	}

	rd, rn, err := armReg2(ops, "cls")
	if err != nil {
		return nil, err
	}

	return Cls{
		rd: rd,
		rn: rn,
	}, nil
}

func newLslReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("lsl: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "lsl")
	if err != nil {
		return nil, err
	}

	return LslReg{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newLsrReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("lsr: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "lsr")
	if err != nil {
		return nil, err
	}

	return LsrReg{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newAsrReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("asr: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "asr")
	if err != nil {
		return nil, err
	}

	return AsrReg{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newRorReg(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("ror: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "ror")
	if err != nil {
		return nil, err
	}

	return RorReg{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newUdiv(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("udiv: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "udiv")
	if err != nil {
		return nil, err
	}

	return Udiv{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newSdiv(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("sdiv: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "sdiv")
	if err != nil {
		return nil, err
	}

	return Sdiv{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newSmulh(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("smulh: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "smulh")
	if err != nil {
		return nil, err
	}

	return Smulh{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newUmulh(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("umulh: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "umulh")
	if err != nil {
		return nil, err
	}

	return Umulh{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newAdc(ops []vOp) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("adc: want rd, rn, rm")
	}

	rd, rn, rm, err := armReg3(ops, "adc")
	if err != nil {
		return nil, err
	}

	return Adc{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

func newCcmp(ops []vOp) (Instr, error) {
	if len(ops) != 4 || ops[2].kind != armOpImm || ops[3].sym == "" {
		return nil, errors.New("ccmp: want rn, rm, #imm, cond")
	}

	rn, rm, err := armReg2(ops, "ccmp")
	if err != nil {
		return nil, err
	}

	v := ops[2].num
	if v < 0 || v > 15 {
		return nil, errors.New("ccmp: bad imm")
	}

	return Ccmp{
		rn:     rn,
		rm:     rm,
		immVal: uint32(v),
		cond:   ops[3].sym,
	}, nil
}

func newExtr(ops []vOp) (Instr, error) {
	if len(ops) != 4 || ops[3].kind != armOpImm {
		return nil, errors.New("extr: want rd, rn, rm, #lsb")
	}

	rd, rn, rm, err := armReg3(ops, "extr")
	if err != nil {
		return nil, err
	}

	v := ops[3].num
	if v < 0 || v > 63 {
		return nil, errors.New("extr: bad lsb")
	}

	return Extr{
		rd:  rd,
		rn:  rn,
		rm:  rm,
		lsb: uint32(v),
	}, nil
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

		return Msub{Madd: Madd{
			rd: rd,
			rn: rn,
			rm: rm,
			ra: ra,
		}}, nil
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

	return Msub{Madd: Madd{
		rd: rd,
		rn: rn,
		rm: rm,
		ra: zeroReg(rd),
	}}, nil
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

		return Madd{
			rd: rd,
			rn: rn,
			rm: rm,
			ra: ra,
		}, nil
	}

	if len(ops) != 3 {
		return nil, fmt.Errorf("%s: want rd, rn, rm[, ra]", name)
	}

	rd, rn, rm, err := armReg3(ops, name)
	if err != nil {
		return nil, err
	}

	return Madd{
		rd: rd,
		rn: rn,
		rm: rm,
		ra: zeroReg(rd),
	}, nil
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

	return Csel{
		rd:   rd,
		rn:   rn,
		rm:   rm,
		cond: cond,
	}, nil
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
	c := newCsel(rd, rn, rm, cond)
	switch base {
	case "csinc":
		return Csinc{Csel: c}, nil
	case "csinv":
		return Csinv{Csel: c}, nil
	case "csneg":
		return Csneg{Csel: c}, nil
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
