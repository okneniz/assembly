package arm64

// Mov wide assembler constructors: movz/movn/movk rd, #imm16
// [, lsl #N]. movn — arithmetic inversion of imm16 at encoding time.

import "fmt"

// movWideOps — rd, #imm16[, lsl #N] → rd, imm16, hw.
func movWideOps(ops []vOp, name string) (string, uint32, uint32, error) {
	if len(ops) != 2 && (len(ops) != 3 || ops[2].kind != armOpShift) {
		return "", 0, 0, fmt.Errorf("%s: want rd, #imm16[, lsl #N]", name)
	}

	rd, err := wantAReg(ops[0], name)
	if err != nil {
		return "", 0, 0, err
	}

	if ops[1].kind != armOpImm || ops[1].sym != "" {
		return "", 0, 0, fmt.Errorf("%s: immediate expected", name)
	}

	v := ops[1].num
	if v < 0 || v > 0xffff {
		return "", 0, 0, fmt.Errorf("%s: imm16 expected", name)
	}

	hw := int64(0)
	if len(ops) == 3 {
		sh := shiftAmt(ops[2])
		if sh%16 != 0 || sh > 48 {
			return "", 0, 0, fmt.Errorf("%s: bad shift %d", name, sh)
		}

		hw = sh / 16
	}

	return rd, uint32(v), uint32(hw), nil
}

func newMovz(ops []vOp) (Instr, error) {
	rd, imm, hw, err := movWideOps(ops, "movz")
	if err != nil {
		return nil, err
	}

	return Movz{
		rd:    rd,
		imm16: imm,
		hw:    hw,
	}, nil
}

func newMovn(ops []vOp) (Instr, error) {
	rd, imm, hw, err := movWideOps(ops, "movn")
	if err != nil {
		return nil, err
	}

	return Movn{
		rd:    rd,
		imm16: imm,
		hw:    hw,
	}, nil
}

func newMovk(ops []vOp) (Instr, error) {
	rd, imm, hw, err := movWideOps(ops, "movk")
	if err != nil {
		return nil, err
	}

	return Movk{
		rd:    rd,
		imm16: imm,
		hw:    hw,
	}, nil
}
