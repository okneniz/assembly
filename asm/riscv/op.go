package riscv

// Unevaluated operands (the syntax layer model): a register, an
// expression, or memory off(reg). Value slots are asm/expr
// expressions; evaluation happens in resolve (see resolve.go), plain
// numbers go into arch/riscv (arch.Op).

import (
	"errors"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm/expr"
)

// Op is an instruction operand: a register, an expression, or memory.
type Op struct {
	reg string
	e   *expr.Expr
	mem *Mem
}

// Mem is the memory operand off(reg); off == nil means bare "(rs1)"
// (AMO).
type Mem struct {
	base string
	off  *expr.Expr
}

// OpReg returns a register by its canonical name.
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

// OpMem returns memory with a numeric offset (pseudo expansion).
func OpMem(base string, off int64) Op {
	return Op{mem: &Mem{base: base, off: expr.Num(off)}}
}

// OpMemExpr returns memory with an expression offset (grammar;
// off == nil means bare "(rs1)").
func OpMemExpr(base string, off *expr.Expr) Op {
	return Op{mem: &Mem{base: base, off: off}}
}

// --- operand access (for pseudo) ---

// Reg returns the register name ("" if the operand is not a register).
func (o Op) Reg() string {
	return o.reg
}

// Expr returns the operand expression (nil if not an expression).
func (o Op) Expr() *expr.Expr {
	return o.e
}

// Mem returns the memory operand (nil if not memory).
func (o Op) Mem() *Mem {
	return o.mem
}

// IsReg reports whether the operand is a register.
func (o Op) IsReg() bool {
	return o.reg != ""
}

// IsMem reports whether the operand is memory.
func (o Op) IsMem() bool {
	return o.mem != nil
}

// BaseReg returns the base register of the memory operand.
func (m *Mem) BaseReg() string {
	return m.base
}

// Off returns the offset expression (nil means bare "(rs1)").
func (m *Mem) Off() *expr.Expr {
	return m.off
}

// WantReg validates a register operand (fp means a floating one is
// expected); the name is checked against the arch table (number +
// width).
func WantReg(op Op, fp bool) (string, error) {
	if op.reg == "" {
		return "", errors.New("want register operand")
	}

	if _, err := arch.RegNumOf(op.reg, fp); err != nil {
		return "", err
	}

	return op.reg, nil
}

// WantExpr validates an expression operand (a pseudo's symbolic
// target).
func WantExpr(op Op) (*expr.Expr, error) {
	if op.e == nil {
		return nil, errors.New("want immediate operand")
	}

	return op.e, nil
}
