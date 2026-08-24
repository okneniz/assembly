package riscv

// Evaluation of an unevaluated instruction: expression slots ->
// numbers, keywords (FP rounding modes, CSR names) -> numbers, then
// building the evaluated instruction with arch.BuildInstr. The "slot
// was symbolic" flag (exprNumeric) is absorbed into EncOpts.NoRVC -
// symbolic targets are not compressed so that the assembler's pass
// sizes converge (core contract: Resolve is stable under a
// placeholder environment).

import (
	"fmt"
	"io"

	arch "github.com/okneniz/assembly/arch/riscv"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// exprNumeric reports that the expression contains no symbols
// (structurally). The RVC-compression condition: the decision must be
// the same in both passes, regardless of whether the symbols are
// already resolved.
func exprNumeric(e *expr.Expr) bool {
	if e == nil {
		return true
	}

	switch e.Kind {
	case expr.ExprNum:
		return true
	case expr.ExprSym:
		return false
	case expr.ExprUnary:
		return exprNumeric(e.X)
	case expr.ExprBinary:
		return exprNumeric(e.X) && exprNumeric(e.Y)
	}

	return false
}

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

// resolved is an evaluated instruction: the arch structure + the
// captured address and modes (pure encoding without an environment).
type resolved struct {
	in   arch.Instr
	pc   uint64
	opts arch.EncOpts
}

func (r resolved) Encode(w io.Writer) (int64, error) {
	return r.in.Encode(w, r.pc, r.opts)
}

// resolve evaluates the expression slots via ctx and builds the
// evaluated instruction. optNoRVC is the .option mode at the
// instruction's moment; symbolic (a structurally symbolic slot) is
// added to it.
func (in instr) resolve(ctx asm.Ctx, optNoRVC bool) (asm.Resolved, error) {
	ops := make([]arch.Op, 0, len(in.ops))
	symbolic := false

	for idx, o := range in.ops {
		switch {
		case o.reg != "":
			ops = append(ops, arch.OpReg(o.reg))
		case o.mem != nil:
			v, num, err := evalExpr(o.mem.off, ctx)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", in.mnem, err)
			}

			symbolic = symbolic || !num
			ops = append(ops, arch.OpMem(o.mem.base, v))
		default:
			if v, ok := keywordNum(in.mnem, idx, o.e); ok {
				ops = append(ops, arch.OpNum(v))
				continue
			}

			v, num, err := evalExpr(o.e, ctx)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", in.mnem, err)
			}

			symbolic = symbolic || !num
			ops = append(ops, arch.OpNum(v))
		}
	}

	st, err := arch.BuildInstr(in.mnem, ops)
	if err != nil {
		return nil, err
	}

	return resolved{
		in:   st,
		pc:   ctx.Addr(),
		opts: arch.EncOpts{NoRVC: optNoRVC || symbolic},
	}, nil
}

// evalExpr evaluates an expression slot (nil means bare memory
// "(rs1)": zero); num reports that the operand is structurally
// numeric (without symbols).
func evalExpr(e *expr.Expr, ctx asm.Ctx) (int64, bool, error) {
	if e == nil {
		return 0, true, nil
	}

	v, err := e.Eval(ctx.Resolve)
	if err != nil {
		return 0, false, err
	}

	return v, exprNumeric(e), nil
}
