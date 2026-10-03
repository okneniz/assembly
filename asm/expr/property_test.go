package expr_test

// Property tests (oh-snap) of GAS expressions: the round trip tree → text →
// tree through the canonical render, literals in all formats by base,
// parser robustness on junk strings, and evaluator robustness on trees.
// Generators are arb/expr; the seed is ASSEMBLY_SEED, as in the root suite
// (property_test.go).

import (
	"iter"
	"math"
	mrnd "math/rand/v2"
	"os"
	"slices"
	"strconv"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	arbx "github.com/okneniz/assembly/arb/expr"
	"github.com/okneniz/assembly/asm/expr"
)

// seedRnd is a deterministic generator: the seed comes from ASSEMBLY_SEED
// (default 42), logged to reproduce a failure (a copy of the root suite
// helper).
func seedRnd(t *testing.T) *mrnd.Rand {
	t.Helper()

	seed := uint64(42)
	if s := os.Getenv("ASSEMBLY_SEED"); s != "" {
		if v, err := strconv.ParseUint(s, 0, 64); err == nil {
			seed = v
		}
	}

	t.Logf("seed: %d (ASSEMBLY_SEED)", seed)
	return arb.Rnd(seed)
}

// renderExpr is the canonical text of a tree: every internal node in
// parentheses, a binary operation with spaces around it ("(1 / 2)" - never
// "//"), a unary one fused with the operand ("(-x)"). Numbers are decimal
// without leading zeros; the generator leaves are non-negative, so the
// render is injective on the image of the grammar (a negative literal "-5"
// parses as unary minus).
func renderExpr(e *expr.Expr) string {
	switch e.Kind {
	case expr.ExprNum:
		return strconv.FormatInt(e.Num, 10)
	case expr.ExprSym:
		return e.Sym
	case expr.ExprUnary:
		return "(" + e.Op + renderExpr(e.X) + ")"
	case expr.ExprBinary:
		return "(" + renderExpr(e.X) + " " + e.Op + " " + renderExpr(e.Y) + ")"
	}

	return "?"
}

// renderBare is the same tree with BARE unary chains: a sign run over a
// leaf or another unary fuses without parentheses ("+-5", "-~x") - the
// spelling an operand takes in a source line. A unary over a BINARY node
// keeps its parentheses ("+(1 + 2)": without them the chain would
// re-associate into binary(+, unary(+, 1), 2)).
func renderBare(e *expr.Expr) string {
	switch e.Kind {
	case expr.ExprNum:
		return strconv.FormatInt(e.Num, 10)
	case expr.ExprSym:
		return e.Sym
	case expr.ExprUnary:
		if e.X.Kind == expr.ExprBinary {
			return "(" + e.Op + renderBare(e.X) + ")"
		}

		return e.Op + renderBare(e.X)
	case expr.ExprBinary:
		return "(" + renderBare(e.X) + " " + e.Op + " " + renderBare(e.Y) + ")"
	}

	return "?"
}

// TestPropertyParseRenderRoundTrip is the "round trip" property: the
// canonical text of a tree parses back into the same tree. Equality is by
// ExprKey: the key is injective on structure and it is the same predicate
// as literal pool deduplication. Both renders walk: the parenthesized
// canonical form and the bare unary-chain spelling.
func TestPropertyParseRenderRoundTrip(t *testing.T) {
	ohsnap.Check(t, 100000, arbx.Tree(seedRnd(t)), func(e *expr.Expr) bool {
		for _, src := range []string{renderExpr(e), renderBare(e)} {
			back, err := expr.ParseExpr(src)
			if err != nil {
				t.Logf("%q: parse: %v", src, err)
				return false
			}

			if expr.ExprKey(back) != expr.ExprKey(e) {
				t.Logf("%q: %s ≠ %s", src, expr.ExprKey(back), expr.ExprKey(e))
				return false
			}
		}

		return true
	})
}

// fmtRow is a row of the literal formats table: the render of a value and
// its domain. The domain is the range of the literal itself: hex/bin are
// read as uint64 (int64 overflow wraps), oct/dec use ParseInt, the signed
// range; char is printable ASCII.
type fmtRow struct {
	name   string
	render func(uint64) string
	from   uint64
	to     uint64
}

// charLit is a character literal of a printable rune; the quote and the
// backslash are escaped (a bare quote would close the literal).
func charLit(v uint64) string {
	switch r := rune(v); r {
	case '\'':
		return `'\''`
	case '\\':
		return `'\\'`
	default:
		return "'" + string(r) + "'"
	}
}

// TestPropertyLiteralFormats is the "formats" property: a literal rendered
// by base parses and evaluates to the written value (Eval == int64(v) -
// the wrap for hex/bin beyond int64 is documented by the domain).
func TestPropertyLiteralFormats(t *testing.T) {
	rows := []fmtRow{
		{
			"hex",
			func(v uint64) string { return "0x" + strconv.FormatUint(v, 16) },
			0,
			math.MaxUint64,
		},
		{
			"bin",
			func(v uint64) string { return "0b" + strconv.FormatUint(v, 2) },
			0,
			math.MaxUint64,
		},
		{"oct", func(v uint64) string { return "0" + strconv.FormatUint(v, 8) }, 0, math.MaxInt64},
		{"dec", func(v uint64) string { return strconv.FormatUint(v, 10) }, 0, math.MaxInt64},
		{"char", charLit, 32, 126},
	}

	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			ohsnap.Check(
				t,
				100000,
				ohsnap.ArbitraryUint64(seedRnd(t), row.from, row.to),
				func(v uint64) bool {
					src := row.render(v)

					e, err := expr.ParseExpr(src)
					if err != nil {
						t.Logf("%q: parse: %v", src, err)
						return false
					}

					got, err := e.Eval(nil)
					if err != nil {
						t.Logf("%q: eval: %v", src, err)
						return false
					}

					if got != int64(v) {
						t.Logf("%q: %d ≠ %d", src, got, int64(v))
						return false
					}

					return true
				},
			)
		})
	}

	// Escapes are a fixed table (non-printable runes are not generated by
	// range); pinning data, not a generator.
	t.Run("char-escape", func(t *testing.T) {
		escapes := []struct {
			src  string
			want int64
		}{
			{`'\n'`, '\n'},
			{`'\t'`, '\t'},
			{`'\r'`, '\r'},
			{`'\0'`, 0},
			{`'\\'`, '\\'},
			{`'\''`, '\''},
			{`'\"'`, '"'},
		}
		for _, esc := range escapes {
			e, err := expr.ParseExpr(esc.src)
			require.NoError(t, err, "ParseExpr(%q)", esc.src)

			v, err := e.Eval(nil)
			require.NoError(t, err, "Eval(%q)", esc.src)
			require.Equal(t, esc.want, v, "Eval(%q)", esc.src)
		}
	})
}

// TestPropertyParseRobustness is the "robustness" property: a junk string
// over the lexicon runes either parses or is rejected with an error; a
// panic or a hang would fail the test. The outcome does not matter - the
// totality of ParseExpr does.
func TestPropertyParseRobustness(t *testing.T) {
	ohsnap.Check(t, 100000, arbx.Junk(seedRnd(t)), func(s string) bool {
		_, err := expr.ParseExpr(s)
		_ = err // the outcome does not matter: the property is totality (no panic), not the parse outcome
		return true
	})
}

// propResolve is a total resolver: every symbol → 42 (symbol values are not
// the subject of the property; totality is).
func propResolve(string) (uint64, bool) {
	return 42, true
}

// TestPropertyEvalRobustness is the "evaluator robustness" property: Eval is
// total on any tree - either a value or an error (shift outside 0..63,
// division by zero), but not a panic.
func TestPropertyEvalRobustness(t *testing.T) {
	ohsnap.Check(t, 100000, arbx.Tree(seedRnd(t)), func(e *expr.Expr) bool {
		_, err := e.Eval(propResolve)
		_ = err // an error is allowed (shift/division by zero), a panic is not
		return true
	})
}

// precLevels - the binary ladder from the loosest to the tightest
// level. The table is the SPEC of the binding order (the grammar's
// binLevel chain): a grammar that deviates from it fails the property.
var precLevels = [][]string{
	{"|"},
	{"^"},
	{"&"},
	{"<<", ">>"},
	{"+", "-"},
	{"*", "/", "%"},
}

// precOf - the level index of a binary operator (the tighter, the
// bigger).
func precOf(op string) int {
	for l, ops := range precLevels {
		if slices.Contains(ops, op) {
			return l
		}
	}

	panic("unknown operator " + op)
}

// precAllBinOps - the binary operators in ladder order.
func precAllBinOps() []string {
	var out []string
	for _, ops := range precLevels {
		out = append(out, ops...)
	}

	return out
}

// precUnOps - the unary operators (the tightest level, above every
// binary one).
var precUnOps = []string{"-", "~", "+"}

// bin - a binary node.
func bin(op string, x, y *expr.Expr) *expr.Expr {
	return expr.NewExpr(expr.ExprBinary, 0, "", op, x, y)
}

// un - a unary node.
func un(op string, x *expr.Expr) *expr.Expr {
	return expr.NewExpr(expr.ExprUnary, 0, "", op, x, nil)
}

// precLeaf - a leaf of a precedence text: its spelling and its tree.
type precLeaf struct {
	text string
	tree *expr.Expr
}

// precLeaves - three random leaves (dec/hex/bin/oct literals and
// symbols - the spellings that can neighbor an operator in a source
// line). Shrinking is structural only (the leaf axis carries no law).
type precLeaves struct {
	a, b, c precLeaf
}

type precLeafGen struct {
	rnd *mrnd.Rand
}

func (g precLeafGen) leaf() precLeaf {
	v := int64(g.rnd.IntN(1 << 12))
	switch g.rnd.IntN(5) {
	case 0:
		return precLeaf{strconv.FormatInt(v, 10), expr.Num(v)}
	case 1:
		return precLeaf{"0x" + strconv.FormatInt(v, 16), expr.Num(v)}
	case 2:
		return precLeaf{"0b" + strconv.FormatInt(v&0xff, 2), expr.Num(v & 0xff)}
	case 3:
		return precLeaf{"0" + strconv.FormatInt(v&0x3ff, 8), expr.Num(v & 0x3ff)}
	default:
		s := []string{"a", "zz", "q7", "sym_2", ".L"}[g.rnd.IntN(5)]
		return precLeaf{s, expr.Sym(s)}
	}
}

func (g precLeafGen) Generate() iter.Seq[precLeaves] {
	return arb.Stream(func() precLeaves {
		return precLeaves{
			a: g.leaf(),
			b: g.leaf(),
			c: g.leaf(),
		}
	})
}

func (precLeafGen) Shrink(precLeaves) iter.Seq[precLeaves] {
	return slices.Values([]precLeaves{})
}

// TestPropertyPrecedence - the binding order of the grammar against the
// spec table: EVERY ordered pair of binary operators (the equal pair is
// the left-associativity law) and every unary x binary combination,
// over random literal/symbol leaves. The check is structural (ExprKey):
// "a x b y c" must parse into the tree the table prescribes.
func TestPropertyPrecedence(t *testing.T) {
	check := func(t *testing.T, src string, want *expr.Expr) bool {
		t.Helper()

		got, err := expr.ParseExpr(src)
		if err != nil {
			t.Logf("%q: parse: %v", src, err)
			return false
		}

		if expr.ExprKey(got) != expr.ExprKey(want) {
			t.Logf("%q: %s ≠ %s", src, expr.ExprKey(got), expr.ExprKey(want))
			return false
		}

		return true
	}

	for _, x := range precAllBinOps() {
		for _, y := range precAllBinOps() {
			t.Run(x+"_"+y, func(t *testing.T) {
				ohsnap.Check(t, 40, precLeafGen{rnd: seedRnd(t)}, func(p precLeaves) bool {
					src := p.a.text + " " + x + " " + p.b.text + " " + y + " " + p.c.text

					// prec(x) >= prec(y): x folds first - either x is
					// tighter, or the level is equal and the parse is
					// left-associative
					if precOf(x) >= precOf(y) {
						return check(t, src, bin(y, bin(x, p.a.tree, p.b.tree), p.c.tree))
					}

					return check(t, src, bin(x, p.a.tree, bin(y, p.b.tree, p.c.tree)))
				})
			})
		}
	}

	for _, u := range precUnOps {
		for _, x := range precAllBinOps() {
			t.Run("un"+u+"_"+x, func(t *testing.T) {
				ohsnap.Check(t, 40, precLeafGen{rnd: seedRnd(t)}, func(p precLeaves) bool {
					// the unary binds tighter than any binary operator,
					// on either side of it
					return check(t, u+p.a.text+" "+x+" "+p.b.text, bin(x, un(u, p.a.tree), p.b.tree)) &&
						check(t, p.a.text+" "+x+" "+u+p.b.text, bin(x, p.a.tree, un(u, p.b.tree)))
				})
			})
		}
	}
}
