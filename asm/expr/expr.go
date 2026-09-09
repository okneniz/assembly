// Package expr is GAS expressions: the AST, a parsec grammar with
// precedences (Chainl1), and an evaluator. Precedences (from low to high):
// | ^ & << >> + - * / %; unary - ~ +. Symbols are evaluated through the
// assembler's resolver (labels, .set constants, "."); the package itself
// knows nothing about symbols.
package expr

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"
)

// ExprKind is the AST node type of an expression.
type ExprKind int

const (
	ExprNum    ExprKind = iota // numeric literal
	ExprSym                    // symbol (label, .set constant, ".")
	ExprUnary                  // -x, ~x, +x
	ExprBinary                 // x op y
)

// Expr is an AST node of an expression.
type Expr struct {
	Kind ExprKind
	Num  int64  // ExprNum
	Sym  string // ExprSym
	Op   string // ExprUnary / ExprBinary
	X, Y *Expr  // operands
}

func NewExpr(kind ExprKind, num int64, sym string, op string, x *Expr, y *Expr) *Expr {
	return &Expr{
		Kind: kind,
		Num:  num,
		Sym:  sym,
		Op:   op,
		X:    x,
		Y:    y,
	}
}

// Num builds a numeric node.
func Num(n int64) *Expr {
	return NewExpr(ExprNum, n, "", "", nil, nil)
}

// Sym builds a node referencing a symbol.
func Sym(s string) *Expr {
	return NewExpr(ExprSym, 0, s, "", nil, nil)
}

// Eval evaluates the expression; resolve supplies symbol values (may be
// nil).
func (e *Expr) Eval(resolve func(string) (uint64, bool)) (int64, error) {
	if e == nil {
		return 0, errors.New("empty expression")
	}

	switch e.Kind {
	case ExprNum:
		return e.Num, nil
	case ExprSym:
		if resolve == nil {
			return 0, fmt.Errorf("undefined symbol %q", e.Sym)
		}

		v, ok := resolve(e.Sym)
		if !ok {
			return 0, fmt.Errorf("undefined symbol %q", e.Sym)
		}

		return int64(v), nil
	case ExprUnary:
		v, err := e.X.Eval(resolve)
		if err != nil {
			return 0, err
		}

		switch e.Op {
		case "-":
			return -v, nil
		case "~":
			return ^v, nil
		case "+":
			return v, nil
		}

		return 0, fmt.Errorf("bad unary operator %q", e.Op)
	case ExprBinary:
		x, err := e.X.Eval(resolve)
		if err != nil {
			return 0, err
		}

		y, err := e.Y.Eval(resolve)
		if err != nil {
			return 0, err
		}

		switch e.Op {
		case "|":
			return x | y, nil
		case "^":
			return x ^ y, nil
		case "&":
			return x & y, nil
		case "<<":
			if y < 0 || y > 63 {
				return 0, fmt.Errorf("shift by %d", y)
			}

			return x << uint(y), nil
		case ">>":
			if y < 0 || y > 63 {
				return 0, fmt.Errorf("shift by %d", y)
			}

			return x >> uint(y), nil
		case "+":
			return x + y, nil
		case "-":
			return x - y, nil
		case "*":
			return x * y, nil
		case "/":
			if y == 0 {
				return 0, errors.New("division by zero")
			}

			return x / y, nil
		case "%":
			if y == 0 {
				return 0, errors.New("division by zero")
			}

			return x % y, nil
		}

		return 0, fmt.Errorf("bad binary operator %q", e.Op)
	}

	return 0, fmt.Errorf("bad expression kind %d", e.Kind)
}

// --- expression grammar ---

// op2 is a binary operator combinator for sym, building an ExprBinary node.
// Special case "/": a single slash is division, but "//" starts a comment
// (arm/riscv Comment), such an operator does not match (Try restores the
// position).
func op2(sym string) parsec.Combinator[rune, parsecstrings.Position, parsec.BinaryOp[*Expr]] {
	str := parsecstrings.String("operator "+sym, sym)
	isDiv := sym == "/"

	scan := func(buf parsec.Buffer[rune, parsecstrings.Position]) (string, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		got, err := str(buf)
		if err != nil {
			return "", err
		}

		if isDiv {
			if r, ok := peekRune(buf); ok && r == '/' {
				return "", parsec.NewParseError(pos, "comment start //")
			}
		}

		return got, nil
	}

	return parsecstrings.Cast(
		parsecstrings.Try(scan),
		func(string) (parsec.BinaryOp[*Expr], error) {
			return func(x, y *Expr) *Expr {
				return NewExpr(ExprBinary, 0, "", sym, x, y)
			}, nil
		},
	)
}

// binLevel is a precedence level: strictChainl1(term, op1 | op2).
func binLevel(term Combinator,
	ops ...parsec.Combinator[rune, parsecstrings.Position, parsec.BinaryOp[*Expr]],
) Combinator {
	return strictChainl1(term, parsecstrings.Choice("operator", ops...))
}

// strictChainl1 is like parsec.Chainl1, but an operator without a right
// operand is an error, not silent absorption ("1+" must not assemble as
// "1").
func strictChainl1(
	term Combinator,
	op parsec.Combinator[rune, parsecstrings.Position, parsec.BinaryOp[*Expr]],
) Combinator {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (*Expr, parsec.Error[parsecstrings.Position]) {
		x, err := term(buf)
		if err != nil {
			return nil, err
		}

		rest := x
		for !buf.IsEOF() {
			save := buf.Position()
			skipWS(buf)
			f, oerr := op(buf)
			if oerr != nil {
				if rerr := Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				break
			}

			y, yerr := term(buf)
			if yerr != nil {
				return nil, parsec.NewParseError(save, "operator without right operand", yerr)
			}

			rest = f(rest, y)
		}

		return rest, nil
	}
}

// grammar is the precedence ladder of the expression language: primary
// (parenthesized expression, character literal, local reference, number,
// symbol), unary - ~ +, then the binary levels from | (lowest) to * / %
// (highest). It is a STRUCT because the grammar is mutually recursive:
// primary -> "(expr)" -> expr -> ... -> primary; the inner closures
// reference the fields, which are all assigned before the first parse.
//
// Precedences (low → high): | ^ & << >> + - * / %.
type grammar struct {
	primary Combinator
	unary   Combinator
	mul     Combinator
	add     Combinator
	shift   Combinator
	and     Combinator
	xor     Combinator
	or      Combinator
	expr    Combinator
}

// newGrammar builds the whole ladder once; the ready Combinators are then
// reused (captured by the consumers' grammars).
func makeGrammar() *grammar {
	g := &grammar{}

	lparen := parsecstrings.Try(parsecstrings.Eq("'('", '('))
	rparen := parsecstrings.Try(parsecstrings.Eq("')'", ')'))

	// the parenthesized body refers to the top level lazily (through the
	// field): the ladder above it is assigned below
	parenBody := func(buf parsec.Buffer[rune, parsecstrings.Position]) (*Expr, parsec.Error[parsecstrings.Position]) {
		return g.expr(buf)
	}

	g.primary = parsecstrings.Choice(
		"operand",
		parsecstrings.Try(parsecstrings.Between(
			lparen,
			parenBody,
			parsecstrings.SkipMany(MakeSpaceParser(), rparen),
		)),
		parsecstrings.Try(makeCharLitParser()),
		parsecstrings.Try(
			makeLocalRefParser(),
		), // "0b"/"1f" before the number: "0b1010" is cut off by the lookahead in newLocalRef
		parsecstrings.Try(makeNumberParser()),
		parsecstrings.Try(makeSymbolParser()), // identifier; "." and ".+8" give Sym(".")/binary
	)

	minus := parsecstrings.Try(parsecstrings.Eq("'-'", '-'))
	tilde := parsecstrings.Try(parsecstrings.Eq("'~'", '~'))
	plus := parsecstrings.Try(parsecstrings.Eq("'+'", '+'))

	g.unary = func(buf parsec.Buffer[rune, parsecstrings.Position]) (*Expr, parsec.Error[parsecstrings.Position]) {
		skipWS(buf)
		for _, u := range []struct {
			op string
			c  parsec.Combinator[rune, parsecstrings.Position, rune]
		}{
			{op: "-", c: minus},
			{op: "~", c: tilde},
			{op: "+", c: plus},
		} {
			if _, err := u.c(buf); err != nil {
				continue
			}

			x, xerr := g.unary(buf)
			if xerr != nil {
				return nil, xerr
			}

			return NewExpr(ExprUnary, 0, "", u.op, x, nil), nil
		}

		return g.primary(buf)
	}

	g.mul = binLevel(g.unary, op2("*"), op2("/"), op2("%"))
	g.add = binLevel(g.mul, op2("+"), op2("-"))
	g.shift = binLevel(g.add, op2("<<"), op2(">>"))
	g.and = binLevel(g.shift, op2("&"))
	g.xor = binLevel(g.and, op2("^"))
	g.or = binLevel(g.xor, op2("|"))
	g.expr = g.or

	return g
}

// CExpr is the whole GAS expression grammar (numbers/symbols/operators):
// a fresh ladder per call; consumers capture the returned value once and
// reuse it (it holds no mutable state).
func MakeExprParser() Combinator {
	return makeGrammar().expr
}

// ParseExpr parses an expression from a string (for tests and utilities);
// the whole text must be consumed.
func ParseExpr(s string) (*Expr, error) {
	cExpr := MakeExprParser()

	body := func(buf parsec.Buffer[rune, parsecstrings.Position]) (*Expr, parsec.Error[parsecstrings.Position]) {
		e, err := cExpr(buf)
		if err != nil {
			return nil, err
		}

		if !buf.IsEOF() {
			return nil, parsec.NewParseError(buf.Position(), "unexpected trailing characters")
		}

		return e, nil
	}
	e, perr := parsecstrings.ParseString(s, body)
	if perr != nil {
		return nil, fmt.Errorf("%w (at %s)", perr, perr.Position())
	}

	return e, nil
}

// ExprKey is the canonical string of an expression (for literal pool
// deduplication: identical literals give the same key given the same
// structure).
func ExprKey(e *Expr) string {
	switch e.Kind {
	case ExprNum:
		return strconv.FormatInt(e.Num, 10)
	case ExprSym:
		return "s:" + e.Sym
	case ExprUnary:
		return "(" + e.Op + " " + ExprKey(e.X) + ")"
	case ExprBinary:
		return "(" + ExprKey(e.X) + " " + e.Op + " " + ExprKey(e.Y) + ")"
	}

	return "?"
}
