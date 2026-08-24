package arm64

// Unevaluated operands (the syntax layer model): value slots are
// asm/expr expressions; evaluation happens in resolve.go (arch receives
// plain numbers: arch.VOp). The shapes mirror arm64.VOp.

import (
	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/asm/expr"
)

type armOpKind int

const (
	armOpReg armOpKind = iota
	armOpImm           // expression (#... or bare)
	armOpLit           // =expression (literal pool, GAS: ldr xN, =literal)
	armOpMem
	armOpList
	armOpShift  // lsl/lsr/asr/ror #imm - modifier of the previous operand
	armOpExtend // uxtw/sxtw…[#imm]
	armOpFloat  // #1.50000000 (fmov imm)
)

type armOp struct {
	kind  armOpKind
	reg   string // register name without the suffix
	arr   string // the ".8b" suffix -> "8b"
	expr  *expr.Expr
	fval  float64
	mem   *armMem
	list  []armListReg
	shift string // lsl/lsr/asr/ror | extension name
}

// newRegOp/newImmOp/newFloatOp/newShiftOp/newMemOp/newExtendOp/newListOp
// are armOp constructors by shape (they mirror the riscv
// OpReg/OpExpr/...).
func newRegOp(reg string) armOp {
	return armOp{
		kind: armOpReg,
		reg:  reg,
	}
}

func newImmOp(e *expr.Expr) armOp {
	return armOp{
		kind: armOpImm,
		expr: e,
	}
}

func newArmOpLit(e *expr.Expr) armOp {
	return armOp{
		kind: armOpLit,
		expr: e,
	}
}

func newFloatOp(v float64) armOp {
	return armOp{
		kind: armOpFloat,
		fval: v,
	}
}

func newShiftOp(shift string, amt *expr.Expr) armOp {
	return armOp{
		kind:  armOpShift,
		shift: shift,
		expr:  amt,
	}
}

func newMemOp(m *armMem) armOp {
	return armOp{
		kind: armOpMem,
		mem:  m,
	}
}

func newExtendOp(ext string) armOp {
	return armOp{
		kind:  armOpExtend,
		shift: ext,
	}
}

func newListOp() armOp {
	return armOp{
		kind: armOpList,
	}
}

type armMem struct {
	base    string
	offExpr *expr.Expr // [x0, #8]
	offReg  string     // [x0, x1]
	opt     string     // lsl/uxtw/sxtw/…
	optAmt  *expr.Expr // lsl #3
	pre     bool       // [x0, #8]!
	post    *expr.Expr // [x0], #8 (post-index imm)
}

// armListReg is a register list element (an alias of arch.ArmListReg:
// one type for both sides).
type armListReg = arch.ArmListReg

func newArmListReg(reg string, arr string) armListReg {
	return arch.NewArmListReg(reg, arr)
}
