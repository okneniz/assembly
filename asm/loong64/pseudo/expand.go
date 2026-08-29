package pseudo

import (
	"fmt"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
	loong64 "github.com/okneniz/assembly/asm/loong64"
)

// form is one base instruction of an expansion: a mnemonic + operands.
type form struct {
	mnem string
	ops  []loong64.Op
}

// expandable lists the single-word pseudo-mnemonics (the expansion
// table below); the li ladders are handled by the same machinery.
var expandable = []string{
	"nop", "move", "not", "ret", "jr",
	"call", "tail",
	"bltz", "bgez", "bgtz", "blez",
	"li.w", "li.d",
}

// zero is the $zero operand of expansions.
func zero() loong64.Op {
	return loong64.OpReg("$zero")
}

// expandPseudo rewrites a single-word pseudo-instruction into base
// forms. The li ladders evaluate their operand via ctx: a symbolic
// operand takes the worst-case chain LENGTH (the content is whatever
// the current pass sees), so the assembler's two passes converge.
func expandPseudo(mnem string, ops []loong64.Op, ctx asm.Ctx) ([]form, error) {
	switch mnem {
	case "nop":
		return []form{{mnem: "andi", ops: []loong64.Op{zero(), zero(), loong64.OpNum(0)}}}, nil
	case "move":
		rd, rj, err := wantRR("move", ops)
		if err != nil {
			return nil, err
		}

		return []form{{mnem: "or", ops: []loong64.Op{rd, rj, zero()}}}, nil
	case "not":
		rd, rj, err := wantRR("not", ops)
		if err != nil {
			return nil, err
		}

		return []form{{mnem: "nor", ops: []loong64.Op{rd, rj, zero()}}}, nil
	case "ret":
		return []form{{mnem: "jirl",
			ops: []loong64.Op{zero(), loong64.OpReg("$ra"), loong64.OpNum(0)}}}, nil
	case "jr":
		rj, err := wantOne(ops, "jr: want rj")
		if err != nil {
			return nil, err
		}

		return []form{{mnem: "jirl", ops: []loong64.Op{zero(), rj, loong64.OpNum(0)}}}, nil
	case "call":
		t, err := wantOne(ops, "call: want target")
		if err != nil {
			return nil, err
		}

		return []form{{mnem: "bl", ops: []loong64.Op{t}}}, nil
	case "tail":
		t, err := wantOne(ops, "tail: want target")
		if err != nil {
			return nil, err
		}

		return []form{{mnem: "b", ops: []loong64.Op{t}}}, nil
	case "bltz", "bgez":
		rj, t, err := wantTwo(ops, mnem)
		if err != nil {
			return nil, err
		}

		base := "blt"
		if mnem == "bgez" {
			base = "bge"
		}

		return []form{{mnem: base, ops: []loong64.Op{rj, zero(), t}}}, nil
	case "bgtz", "blez":
		rd, t, err := wantTwo(ops, mnem)
		if err != nil {
			return nil, err
		}

		base := "blt"
		if mnem == "blez" {
			base = "bge"
		}

		return []form{{mnem: base, ops: []loong64.Op{zero(), rd, t}}}, nil
	case "li.w":
		return expandLi(ops, false, ctx)
	case "li.d":
		return expandLi(ops, true, ctx)
	}

	return nil, fmt.Errorf("unknown pseudo-instruction %q", mnem)
}

// expandLi builds the immediate-load ladder: a numeric operand takes
// the shortest exact chain (the llvm ladder); a symbolic one (an
// address, in practice) is forced to the full 64-bit chain - as GAS,
// li.w/li.d of a symbol load the whole address.
func expandLi(ops []loong64.Op, is64 bool, ctx asm.Ctx) ([]form, error) {
	what := "li.w"
	if is64 {
		what = "li.d"
	}

	rd, e, err := wantRegExpr(what, ops)
	if err != nil {
		return nil, err
	}

	v, verr := e.Eval(ctx.Resolve)
	if verr != nil {
		return nil, fmt.Errorf("%s: %w", what, verr)
	}

	if !exprNumeric(e) {
		return liWorst64(rd, v), nil
	}

	if is64 {
		return li64(rd, v)
	}

	return li32(rd, v)
}

// li32 - the 32-bit ladders (exact for v representable as int32).
func li32(rd loong64.Op, v int64) ([]form, error) {
	switch {
	case v >= 0 && v <= 4095:
		return []form{{mnem: "ori", ops: []loong64.Op{rd, zero(), loong64.OpNum(v)}}}, nil
	case v >= -2048 && v <= -1:
		return []form{{mnem: "addi.w", ops: []loong64.Op{rd, zero(), loong64.OpNum(v)}}}, nil
	case int64(int32(v)) != v:
		return nil, fmt.Errorf("li.w: value %#x does not fit 32 bits", v)
	case v&0xfff == 0:
		return []form{{mnem: "lu12i.w", ops: []loong64.Op{rd, loong64.OpNum(v >> 12)}}}, nil
	default:
		hi := (v - (v & 0xfff)) >> 12

		return []form{
			{mnem: "lu12i.w", ops: []loong64.Op{rd, loong64.OpNum(hi)}},
			{mnem: "ori", ops: []loong64.Op{rd, rd, loong64.OpNum(v & 0xfff)}},
		}, nil
	}
}

// li64 - the 64-bit ladders.
func li64(rd loong64.Op, v int64) ([]form, error) {
	switch {
	case v >= 0 && v <= 4095:
		return []form{{mnem: "ori", ops: []loong64.Op{rd, zero(), loong64.OpNum(v)}}}, nil
	case v >= -2048 && v <= -1:
		return []form{{mnem: "addi.w", ops: []loong64.Op{rd, zero(), loong64.OpNum(v)}}}, nil
	case int64(int32(v)) == v && v&0xfff == 0:
		return []form{{mnem: "lu12i.w", ops: []loong64.Op{rd, loong64.OpNum(v >> 12)}}}, nil
	case int64(int32(v)) == v:
		hi := (v - (v & 0xfff)) >> 12

		return []form{
			{mnem: "lu12i.w", ops: []loong64.Op{rd, loong64.OpNum(hi)}},
			{mnem: "ori", ops: []loong64.Op{rd, rd, loong64.OpNum(v & 0xfff)}},
		}, nil
	case v&0xfff == 0:
		return []form{
			{mnem: "lu12i.w", ops: []loong64.Op{rd, loong64.OpNum(sext20((v >> 12) & 0xfffff))}},
			{mnem: "lu32i.d", ops: []loong64.Op{rd, loong64.OpNum(sext20((v >> 32) & 0xfffff))}},
			{mnem: "lu52i.d", ops: []loong64.Op{rd, rd, loong64.OpNum(v >> 52)}},
		}, nil
	default:
		return liWorst64(rd, v), nil
	}
}

// liWorst64 - the fixed 4-word chain of li.d.
func liWorst64(rd loong64.Op, v int64) []form {
	return []form{
		{mnem: "lu12i.w", ops: []loong64.Op{rd, loong64.OpNum(sext20((v >> 12) & 0xfffff))}},
		{mnem: "ori", ops: []loong64.Op{rd, rd, loong64.OpNum(v & 0xfff)}},
		{mnem: "lu32i.d", ops: []loong64.Op{rd, loong64.OpNum(sext20((v >> 32) & 0xfffff))}},
		{mnem: "lu52i.d", ops: []loong64.Op{rd, rd, loong64.OpNum(v >> 52)}},
	}
}

// sext20 folds an unsigned 20-bit field into its signed si20 value.
func sext20(v int64) int64 {
	if v >= 0x80000 {
		return v - 0x100000
	}

	return v
}

// --- operand shorthands of the expansions ---

// wantRR - "rd, rj".
func wantRR(what string, ops []loong64.Op) (loong64.Op, loong64.Op, error) {
	if len(ops) != 2 {
		return loong64.Op{}, loong64.Op{}, fmt.Errorf("%s: want rd, rj", what)
	}

	return ops[0], ops[1], nil
}

// wantOne - a single operand.
func wantOne(ops []loong64.Op, what string) (loong64.Op, error) {
	if len(ops) != 1 {
		return loong64.Op{}, fmt.Errorf("%s", what)
	}

	return ops[0], nil
}

// wantTwo - "reg, target".
func wantTwo(ops []loong64.Op, what string) (loong64.Op, loong64.Op, error) {
	if len(ops) != 2 {
		return loong64.Op{}, loong64.Op{}, fmt.Errorf("%s: want rj, target", what)
	}

	return ops[0], ops[1], nil
}

// wantRegExpr - "rd, imm".
func wantRegExpr(what string, ops []loong64.Op) (loong64.Op, *expr.Expr, error) {
	if len(ops) != 2 || !ops[0].IsReg() || ops[1].Expr() == nil {
		return loong64.Op{}, nil, fmt.Errorf("%s: want rd, imm", what)
	}

	return ops[0], ops[1].Expr(), nil
}

// exprNumeric reports that the expression contains no symbols
// (structurally): the li chain-length decision must be identical in
// both assembler passes.
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
