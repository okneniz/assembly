package arm64

import "errors"

// Branch assembler constructors: b/bl label, cbz/cbnz rt, label,
// br/blr xn, ret{ xn}, nop. Targets arrive computed (resolveOps of the
// syntax layer).

func newB(ops []vOp) (Instr, error) {
	if len(ops) != 1 {
		return nil, errors.New("b: want target")
	}

	t, err := wantTarget(ops[0], "b")
	if err != nil {
		return nil, err
	}

	return B{target: t}, nil
}

func newBl(ops []vOp) (Instr, error) {
	if len(ops) != 1 {
		return nil, errors.New("bl: want target")
	}

	t, err := wantTarget(ops[0], "bl")
	if err != nil {
		return nil, err
	}

	return Bl{target: t}, nil
}

func newCbz(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("cbz: want rt, target")
	}

	rt, err := wantAReg(ops[0], "cbz")
	if err != nil {
		return nil, err
	}

	t, err := wantTarget(ops[1], "cbz")
	if err != nil {
		return nil, err
	}

	return Cbz{
		rt:     rt,
		target: t,
	}, nil
}

func newCbnz(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("cbnz: want rt, target")
	}

	rt, err := wantAReg(ops[0], "cbnz")
	if err != nil {
		return nil, err
	}

	t, err := wantTarget(ops[1], "cbnz")
	if err != nil {
		return nil, err
	}

	return Cbnz{
		rt:     rt,
		target: t,
	}, nil
}

func newNop(ops []vOp) (Instr, error) {
	if len(ops) != 0 {
		return nil, errors.New("nop expects no operands")
	}

	return Nop{}, nil
}

func newRet(ops []vOp) (Instr, error) {
	if len(ops) > 1 {
		return nil, errors.New("ret expects at most one operand")
	}

	rn := "x30"
	if len(ops) == 1 {
		r, err := wantAReg(ops[0], "ret")
		if err != nil {
			return nil, err
		}

		rn = r
	}

	return Ret{rn: rn}, nil
}

func newBr(ops []vOp) (Instr, error) {
	if len(ops) != 1 {
		return nil, errors.New("br: want xn")
	}

	rn, err := wantAReg(ops[0], "br")
	if err != nil {
		return nil, err
	}

	return Br{rn: rn}, nil
}

func newBlr(ops []vOp) (Instr, error) {
	if len(ops) != 1 {
		return nil, errors.New("blr: want xn")
	}

	rn, err := wantAReg(ops[0], "blr")
	if err != nil {
		return nil, err
	}

	return Blr{rn: rn}, nil
}
