package riscv

import "fmt"

// csrOp - the shared bookkeeping record of a CSR operand: name/number
// (a computed number; the name is recovered from the rvCsrNames table -
// unknown numbers are printed as %#x).
type csrOp struct {
	rd   string
	csr  string // name ("fflags") for a known CSR, otherwise ""
	csrN uint32 // 12-bit number
}

// newCsrOp - from a number (the syntax layer has evaluated the name/expression to a number).
func newCsrOp(rd string, num int64) csrOp {
	return csrOp{
		rd:   rd,
		csr:  rvCsrNames[uint32(num)],
		csrN: uint32(num),
	}
}

// csrBits - the CSR number for encoding.
func (c csrOp) csrBits() uint32 {
	return c.csrN
}

// text - the displayed name (for pseudo checks and JSON fields).
func (c csrOp) text() string {
	if c.csr != "" {
		return c.csr
	}

	return fmt.Sprintf("%#x", c.csrN)
}

// zimmBits - the 5-bit zimm for encoding.
func zimmBits(m imm) (uint32, error) {
	if m.val < 0 || m.val > 31 {
		return 0, fmt.Errorf("zimm %d out of 5-bit range", m.val)
	}

	return uint32(m.val), nil
}

// newCsrr — csrrw/csrrs/csrrc rd, csr, rs1.
func newCsrr(ops []Op, name string) (Instr, error) {
	if len(ops) != 3 {
		return nil, fmt.Errorf("%s: want rd, csr, rs1", name)
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	num, err := wantCSR(ops[1])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	rs1, err := wantReg(ops[2], false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	cf := newCsrOp(rd, num)
	switch name {
	case "csrrw":
		return Csrrw{
			csrOp: cf,
			rs1:   rs1,
		}, nil
	case "csrrs":
		return Csrrs{
			csrOp: cf,
			rs1:   rs1,
		}, nil
	default:
		return Csrrc{
			csrOp: cf,
			rs1:   rs1,
		}, nil
	}
}

// newCsrI — csrrwi/csrrsi/csrrci rd, csr, zimm.
func newCsrI(ops []Op, name string) (Instr, error) {
	if len(ops) != 3 {
		return nil, fmt.Errorf("%s: want rd, csr, zimm", name)
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	num, err := wantCSR(ops[1])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	e, err := wantExpr(ops[2])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	cf := newCsrOp(rd, num)
	switch name {
	case "csrrwi":
		return Csrrwi{
			csrOp: cf,
			zimm:  e,
		}, nil
	case "csrrsi":
		return Csrrsi{
			csrOp: cf,
			zimm:  e,
		}, nil
	default:
		return Csrrci{
			csrOp: cf,
			zimm:  e,
		}, nil
	}
}
