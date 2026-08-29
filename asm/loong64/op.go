package loong64

// Unevaluated operands (the syntax layer model): a register or an
// expression. Value slots are asm/expr expressions; evaluation happens
// in resolve (see resolve.go), plain numbers go into arch/loong64
// (arch.Op). LoongArch has no memory operands - loads/stores take flat
// "rd, rj, offset" operands.

import (
	"errors"

	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/asm/expr"
)

// Op is an instruction operand: a register or an expression.
type Op struct {
	reg string
	e   *expr.Expr
}

// OpReg returns a register operand by its canonical $-name.
func OpReg(name string) Op {
	return Op{reg: name}
}

// OpExpr returns an expression operand.
func OpExpr(e *expr.Expr) Op {
	return Op{e: e}
}

// OpNum returns a numeric operand.
func OpNum(v int64) Op {
	return Op{e: expr.Num(v)}
}

// Reg returns the register name ("" if the operand is not a register).
func (o Op) Reg() string {
	return o.reg
}

// Expr returns the operand expression (nil if not an expression).
func (o Op) Expr() *expr.Expr {
	return o.e
}

// IsReg reports whether the operand is a register.
func (o Op) IsReg() bool {
	return o.reg != ""
}

// WantReg validates a register operand; the name is checked against the
// arch table.
func WantReg(op Op) (string, error) {
	if op.reg == "" {
		return "", errors.New("want register operand")
	}

	if _, err := arch.RegNumOf(op.reg); err != nil {
		return "", err
	}

	return op.reg, nil
}

// WantExpr validates an expression operand (a pseudo's symbolic target).
func WantExpr(op Op) (*expr.Expr, error) {
	if ope := op.Expr(); ope != nil {
		return ope, nil
	}

	return nil, errors.New("want immediate operand")
}
