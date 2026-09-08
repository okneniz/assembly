package riscv

// Evaluation of an unevaluated instruction: expression slots ->
// numbers, keywords (FP rounding modes, CSR names) -> numbers, then
// building the evaluated instruction with arch.BuildInstr. RVC
// compression is purely value-driven (a symbolic target compresses when
// the RESOLVED distance fits the c.* range); size convergence between
// the passes is owned by the core layout relaxation (asm.walkLayout).

import (
	"fmt"
	"io"

	arch "github.com/okneniz/assembly/arch/riscv"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// roundingModes maps FP rounding mode names (an optional operand:
// fadd.s fa1, fa2, rne).
var roundingModes = map[string]uint32{
	"rne": 0, "rtz": 1, "rdn": 2, "rup": 3, "rmm": 4, "dyn": 7,
}

// rmSlot is the position of the optional FP arithmetic rounding mode.
var rmSlot = map[string]int{
	"fadd.s": 3, "fsub.s": 3, "fmul.s": 3, "fdiv.s": 3,
	"fadd.d": 3, "fsub.d": 3, "fmul.d": 3, "fdiv.d": 3,
	"fmadd.s": 4, "fmsub.s": 4, "fnmsub.s": 4, "fnmadd.s": 4,
	"fmadd.d": 4, "fmsub.d": 4, "fnmsub.d": 4, "fnmadd.d": 4,
}

// csrSlot lists mnemonics with a CSR operand in position 1.
var csrSlot = map[string]bool{
	"csrrw": true, "csrrs": true, "csrrc": true,
	"csrrwi": true, "csrrsi": true, "csrrci": true,
}

// keywordNum returns the number for a keyword in slot idx of mnemonic
// mnem (a rounding mode, a CSR name); ok=false means not a keyword
// (an ordinary expression: evaluated by the resolver).
func keywordNum(mnem string, idx int, e *expr.Expr) (int64, bool) {
	if e == nil || e.Kind != expr.ExprSym {
		return 0, false
	}

	if rmSlot[mnem] == idx {
		if v, ok := roundingModes[e.Sym]; ok {
			return int64(v), true
		}

		return 0, false
	}

	if csrSlot[mnem] && idx == 1 {
		if v, ok := arch.CSRNumOf(e.Sym); ok {
			return int64(v), true
		}
	}

	return 0, false
}

// branchForms - the mnemonics whose LAST operand is a pc-relative
// target: both a symbolic expression and a numeric one evaluate to an
// absolute address (verified against clang at a nonzero address), so
// the value shifts by -pc before the structure is built (the arch
// canon stores the offset).
var branchForms = map[string]bool{
	"beq": true, "bne": true, "blt": true, "bge": true,
	"bltu": true, "bgeu": true,
	"jal": true,
}

// resolved is an evaluated instruction: the arch structure + the
// modes (pure encoding without an environment).
type resolved struct {
	in   arch.Instr
	opts arch.EncOpts
}

func (r resolved) Encode(w io.Writer) (int64, error) {
	return r.in.Encode(w, r.opts)
}

// resolve evaluates the expression slots via ctx and builds the
// evaluated instruction; optNoRVC is the .option mode at the
// instruction's moment (compression itself is value-driven - see the
// package comment).
func (in instr) resolve(ctx asm.Ctx, optNoRVC bool) (asm.Resolved, error) {
	ops := make([]arch.Op, 0, len(in.ops))

	for idx, o := range in.ops {
		switch {
		case o.reg != "":
			ops = append(ops, arch.OpReg(o.reg))
		case o.mem != nil:
			v, err := evalExpr(o.mem.off, ctx)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", in.mnem, err)
			}

			ops = append(ops, arch.OpMem(o.mem.base, v))
		default:
			if v, ok := keywordNum(in.mnem, idx, o.e); ok {
				ops = append(ops, arch.OpNum(v))
				continue
			}

			v, err := evalExpr(o.e, ctx)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", in.mnem, err)
			}

			if branchForms[in.mnem] && idx == len(in.ops)-1 {
				v -= int64(ctx.Addr())
			}

			ops = append(ops, arch.OpNum(v))
		}
	}

	st, err := arch.BuildInstr(in.mnem, ops)
	if err != nil {
		return nil, err
	}

	return resolved{
		in:   st,
		opts: arch.EncOpts{NoRVC: optNoRVC},
	}, nil
}

// evalExpr evaluates an expression slot (nil means bare memory
// "(rs1)": zero).
func evalExpr(e *expr.Expr, ctx asm.Ctx) (int64, error) {
	if e == nil {
		return 0, nil
	}

	v, err := e.Eval(ctx.Resolve)
	return v, err
}
