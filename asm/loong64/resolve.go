package loong64

// Evaluation of an unevaluated instruction: expression slots ->
// numbers, then building the evaluated instruction with
// arch.BuildInstr. LA64 words are fixed 32-bit with no compression, so
// the core's placeholder-stability contract holds trivially: Resolve
// sizes never depend on symbol values.

import (
	"fmt"
	"io"

	arch "github.com/okneniz/assembly/arch/loong64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// exprNumeric reports that the expression contains no symbols
// (structurally): only symbolic branch operands shift by pc.
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

// branchForms - the mnemonics whose last operand is a pc-relative byte
// offset: a symbolic expression evaluates to its offset (value - pc),
// a numeric one IS the offset (the assembler convention, verified
// against clang at a nonzero address).
var branchForms = map[string]bool{
	"beq": true, "bne": true, "blt": true, "bge": true,
	"bltu": true, "bgeu": true,
	"beqz": true, "bnez": true, "b": true, "bl": true,
}

// resolved is an evaluated instruction: the arch structure + the
// captured address (pure encoding without an environment).
type resolved struct {
	in arch.Instr
	pc uint64
}

func (r resolved) Encode(w io.Writer) (int64, error) {
	return r.in.Encode(w, r.pc)
}

// resolve evaluates the expression slots via ctx and builds the
// evaluated instruction.
func (in instr) resolve(ctx asm.Ctx) (asm.Resolved, error) {
	ops := make([]arch.Op, 0, len(in.ops))
	for idx, o := range in.ops {
		switch {
		case o.reg != "":
			num, err := arch.RegNumOf(o.reg)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", in.mnem, err)
			}

			ops = append(ops, arch.OpReg(uint8(num)))
		default:
			v, err := o.e.Eval(ctx.Resolve)
			if err != nil {
				return nil, fmt.Errorf("%s: operand %d: %w", in.mnem, idx+1, err)
			}

			if branchForms[in.mnem] && idx == len(in.ops)-1 && !exprNumeric(o.e) {
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
		in: st,
		pc: ctx.Addr(),
	}, nil
}
