// Package expr — arbitrary test generators (oh-snap) for GAS expressions:
// structural *expr.Expr trees (numbers/symbols/unary/binary nodes) and junk
// strings over the expression lexicon. The consumer is the property suite
// asm/expr/property_test.go.
package expr

// Expression tree generator: one generator, one type (*expr.Expr),
// structural growth with a depth budget and per-axis shrink (like the arb/arch families).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/asm/expr"
)

// binOps — the binary operators of the grammar (all ten); vals[0] is the
// shrink target of the operator axis (canonicalization, like arb.Enum).
var binOps = []string{"|", "^", "&", "<<", ">>", "+", "-", "*", "/", "%"}

// unOps — unary operators; "-" is the shrink target.
var unOps = []string{"-", "~", "+"}

// identStart/identCont — symbol name runes: exactly the lexer's predicates
// (expr.IsIdentStart/IsIdentCont), so that every generated name
// re-lexes as a Sym.
const (
	identStart = "._$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	identCont  = identStart + "0123456789"
)

// treeGen — tree generator: depth <= 4, leaves are non-negative uint64>>1
// literals (the grammar has no negative literals — unary minus is a
// separate node) and symbol names, internal nodes are unary/binary.
type treeGen struct {
	rnd *rand.Rand
}

// Tree — an arbitrary expression tree.
func Tree(rnd *rand.Rand) ohsnap.Arbitrary[*expr.Expr] {
	return treeGen{rnd: rnd}
}

func (g treeGen) Generate() *expr.Expr {
	return g.grow(4)
}

// Shrink — one axis at a time: structural simplification (unary → operand,
// binary → operand), leaf (the number is halved toward zero, the name is
// trimmed from the end), in depth (a child is replaced by one of its own
// candidates), operator → vals[0].
func (g treeGen) Shrink(e *expr.Expr) []*expr.Expr {
	if e == nil {
		return nil
	}

	var out []*expr.Expr
	switch e.Kind {
	case expr.ExprUnary:
		out = append(out, e.X)
	case expr.ExprBinary:
		out = append(out, e.X, e.Y)
	case expr.ExprNum:
		for _, v := range arb.Halved(e.Num) {
			out = append(out, expr.Num(v))
		}
	case expr.ExprSym:
		for n := len(e.Sym) - 1; n >= 1; n-- {
			out = append(out, expr.Sym(e.Sym[:n]))
		}
	}

	for _, c := range g.Shrink(e.X) {
		out = append(out, expr.NewExpr(e.Kind, e.Num, e.Sym, e.Op, c, e.Y))
	}

	for _, c := range g.Shrink(e.Y) {
		out = append(out, expr.NewExpr(e.Kind, e.Num, e.Sym, e.Op, e.X, c))
	}

	if e.Kind == expr.ExprUnary && e.Op != unOps[0] {
		out = append(out, expr.NewExpr(expr.ExprUnary, 0, "", unOps[0], e.X, nil))
	}

	if e.Kind == expr.ExprBinary && e.Op != binOps[0] {
		out = append(out, expr.NewExpr(expr.ExprBinary, 0, "", binOps[0], e.X, e.Y))
	}

	return out
}

// grow — a subtree of depth <= depth; a 1/3 chance of a leaf and within the budget.
func (g treeGen) grow(depth int) *expr.Expr {
	if depth == 0 || g.rnd.IntN(3) == 0 {
		if g.rnd.IntN(2) == 0 {
			return expr.Num(int64(g.rnd.Uint64() >> 1)) // [0, MaxInt64]
		}

		return expr.Sym(g.ident())
	}

	if g.rnd.IntN(4) == 0 {
		return expr.NewExpr(expr.ExprUnary, 0, "",
			unOps[g.rnd.IntN(len(unOps))], g.grow(depth-1), nil)
	}

	return expr.NewExpr(expr.ExprBinary, 0, "",
		binOps[g.rnd.IntN(len(binOps))], g.grow(depth-1), g.grow(depth-1))
}

// ident — a name [._$a-zA-Z][._$a-zA-Z0-9]{0,6}.
func (g treeGen) ident() string {
	rs := []rune{rune(identStart[g.rnd.IntN(len(identStart))])}
	for n := g.rnd.IntN(7); n > 0; n-- {
		rs = append(rs, rune(identCont[g.rnd.IntN(len(identCont))]))
	}

	return string(rs)
}
