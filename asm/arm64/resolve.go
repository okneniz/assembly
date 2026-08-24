package arm64

// Evaluation of an unevaluated instruction: expression slots -> numbers
// (arch.VOp). Keywords (conditions, sysreg, prfm hints) become names
// (arch.IsKeywordSym); labels/.set become values via the core resolver;
// adr/adrp: a symbolic target in position 1 is an absolute -> offset
// (as the formatter prints); the ldr=literal is the pool slot address
// (the core resolver returns it under the reserved name PoolSelf). The
// core's placeholder pass yields rel=0.

import (
	"errors"

	arch "github.com/okneniz/assembly/arch/arm64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// resolveOps evaluates the operand expression slots via ctx and builds
// the evaluated arch.VOp operands.
func resolveOps(mnem string, ops []armOp, ctx ctx) ([]arch.VOp, error) {
	out := make([]arch.VOp, 0, len(ops))
	for i, op := range ops {
		switch op.kind {
		case armOpReg:
			var idx int64
			lane := false
			if op.expr != nil { // lane suffix index: "v30[1]"
				n, _, err := resolveSlot(op.expr, ctx)
				if err != nil {
					return nil, err
				}

				idx, lane = n, true
			}

			out = append(out, arch.VOpReg(op.reg, op.arr, lane, idx))
		case armOpImm:
			// adr/adrp: a symbolic target in position 1 is an absolute ->
			// offset; a numeric one is already an offset
			if i == 1 && (mnem == "adr" || mnem == "adrp") &&
				op.expr != nil && op.expr.Kind == expr.ExprSym {
				abs, err := op.expr.Eval(ctx.Resolve)
				if err != nil {
					return nil, err
				}

				out = append(out, arch.VOpImm(abs-int64(ctx.Addr), ""))
				continue
			}

			n, sym, err := resolveSlot(op.expr, ctx)
			if err != nil {
				return nil, err
			}

			out = append(out, arch.VOpImm(n, sym))
		case armOpLit:
			// the literal pool slot address (the reserved name PoolSelf);
			// the slot size is the width of rt (ldr xN,=... is 8; wN is 4).
			// The placeholder pass does not know the slot - the direct
			// value suffices (the size is always 4).
			num := int64(0)
			if addr, ok := ctx.Resolve(asm.PoolSelf); ok {
				num = int64(addr)
			} else if v, verr := op.expr.Eval(ctx.Resolve); verr == nil {
				num = v
			}

			out = append(out, arch.VOpLit(num))
		case armOpShift:
			n, _, err := resolveSlot(op.expr, ctx)
			if err != nil {
				return nil, err
			}

			out = append(out, arch.VOpShift(op.shift, n, op.expr != nil))
		case armOpExtend:
			n, _, err := resolveSlot(op.expr, ctx)
			if err != nil {
				return nil, err
			}

			out = append(out, arch.VOpExtend(op.shift, n, op.expr != nil))
		case armOpMem:
			m := op.mem
			vm := arch.NewVMem(m.base, 0, false, m.offReg, m.opt, 0, false, m.pre, 0, false)
			if m.offExpr != nil {
				n, _, err := resolveSlot(m.offExpr, ctx)
				if err != nil {
					return nil, err
				}

				vm = arch.NewVMem(
					m.base,
					n,
					true,
					m.offReg,
					m.opt,
					vm.OptAmt(),
					false,
					m.pre,
					0,
					false,
				)
			}

			if m.optAmt != nil {
				n, _, err := resolveSlot(m.optAmt, ctx)
				if err != nil {
					return nil, err
				}

				vm = arch.NewVMem(
					m.base,
					vm.Off(),
					vm.HasOff(),
					m.offReg,
					m.opt,
					n,
					true,
					m.pre,
					0,
					false,
				)
			}

			if m.post != nil {
				n, _, err := resolveSlot(m.post, ctx)
				if err != nil {
					return nil, err
				}

				vm = arch.NewVMem(
					m.base,
					vm.Off(),
					vm.HasOff(),
					m.offReg,
					m.opt,
					vm.OptAmt(),
					vm.HasOpt(),
					m.pre,
					n,
					true,
				)
			}

			out = append(out, arch.VOpMem(vm))
		case armOpList, armOpFloat:
			// no expression slots
			if op.kind == armOpList {
				out = append(out, arch.VOpList(op.list))
			} else {
				out = append(out, arch.VOpFloat(op.fval))
			}
		default:
			return nil, errors.New("bad operand kind")
		}
	}

	return out, nil
}

// resolveSlot evaluates an expression slot: a bare Sym keyword -> the
// name; everything else -> the value via the resolver (a nil
// expression is zero).
func resolveSlot(e *expr.Expr, ctx ctx) (int64, string, error) {
	if e == nil {
		return 0, "", nil
	}

	if e.Kind == expr.ExprSym && arch.IsKeywordSym(e.Sym) {
		return 0, e.Sym, nil
	}

	v, err := e.Eval(ctx.Resolve)
	if err != nil {
		return 0, "", err
	}

	return v, "", nil
}
