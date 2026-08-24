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

	"github.com/okneniz/parsec/common"
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
func op2(sym string) common.Combinator[rune, parsecstrings.Position, common.BinaryOp[*Expr]] {
	return parsecstrings.Cast(
		parsecstrings.Try(func() common.Combinator[rune, parsecstrings.Position, string] {
			return func(buf common.Buffer[rune, parsecstrings.Position]) (string, common.Error[parsecstrings.Position]) {
				pos := buf.Position()
				got, err := parsecstrings.String("operator "+sym, sym)(buf)
				if err != nil {
					return "", err
				}

				if sym == "/" {
					if r, ok := peekRune(buf); ok && r == '/' {
						return "", common.NewParseError(pos, "comment start //")
					}
				}

				return got, nil
			}
		}()),
		func(string) (common.BinaryOp[*Expr], error) {
			return func(x, y *Expr) *Expr {
				return NewExpr(ExprBinary, 0, "", sym, x, y)
			}, nil
		},
	)
}

// binLevel is a precedence level: strictChainl1(term, op1 | op2).
func binLevel(term common.Combinator[rune, parsecstrings.Position, *Expr],
	ops ...common.Combinator[rune, parsecstrings.Position, common.BinaryOp[*Expr]],
) common.Combinator[rune, parsecstrings.Position, *Expr] {
	return strictChainl1(term, parsecstrings.Choice("operator", ops...))
}

// strictChainl1 is like common.Chainl1, but an operator without a right
// operand is an error, not silent absorption ("1+" must not assemble as
// "1").
func strictChainl1(
	term common.Combinator[rune, parsecstrings.Position, *Expr],
	op common.Combinator[rune, parsecstrings.Position, common.BinaryOp[*Expr]],
) common.Combinator[rune, parsecstrings.Position, *Expr] {
	return func(buf common.Buffer[rune, parsecstrings.Position]) (*Expr, common.Error[parsecstrings.Position]) {
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
				return nil, common.NewParseError(save, "operator without right operand", yerr)
			}

			rest = f(rest, y)
		}

		return rest, nil
	}
}

// cPrimary is a primary expression: (expr) | number | symbol | '.'.
// cUnary is unary -, ~, + (right-associative) | primary.
// Precedence levels (low → high): | ^ & << >> + - * / %.
// Built in init(): the grammar is mutually recursive (primary → parens →
// expression), a var chain would create an initialization cycle.
var (
	cPrimary common.Combinator[rune, parsecstrings.Position, *Expr]
	cUnary   common.Combinator[rune, parsecstrings.Position, *Expr]
	cMul     common.Combinator[rune, parsecstrings.Position, *Expr]
	cAdd     common.Combinator[rune, parsecstrings.Position, *Expr]
	cShift   common.Combinator[rune, parsecstrings.Position, *Expr]
	cAnd     common.Combinator[rune, parsecstrings.Position, *Expr]
	cXor     common.Combinator[rune, parsecstrings.Position, *Expr]
	cOr      common.Combinator[rune, parsecstrings.Position, *Expr]
	cExpr    common.Combinator[rune, parsecstrings.Position, *Expr]
)

func init() {
	cPrimary = parsecstrings.Choice(
		"operand",
		parsecstrings.Try(parsecstrings.Between(
			cLParen,
			func() common.Combinator[rune, parsecstrings.Position, *Expr] {
				return func(buf common.Buffer[rune, parsecstrings.Position]) (*Expr, common.Error[parsecstrings.Position]) {
					return cExpr(buf) // lazy reference for mutual recursion
				}
			}(),
			parsecstrings.SkipMany(cSpace, cRParen),
		)),
		parsecstrings.Try(cCharLit),
		parsecstrings.Try(
			cLocalRef,
		), // "0b"/"1f" before the number: "0b1010" is cut off by the lookahead in cLocalRef
		parsecstrings.Try(cNumber),
		parsecstrings.Try(cSymbolExpr), // identifier; "." and ".+8" give Sym(".")/binary
	)

	cUnary = func(buf common.Buffer[rune, parsecstrings.Position]) (*Expr, common.Error[parsecstrings.Position]) {
		skipWS(buf)
		for _, u := range []struct {
			op string
			c  common.Combinator[rune, parsecstrings.Position, rune]
		}{{
			"-",
			cMinus,
		}, {
			"~",
			cTilde,
		}, {
			"+",
			cPlus,
		}} {
			if _, err := u.c(buf); err != nil {
				continue
			}

			x, xerr := cUnary(buf)
			if xerr != nil {
				return nil, xerr
			}

			return NewExpr(ExprUnary, 0, "", u.op, x, nil), nil
		}

		return cPrimary(buf)
	}

	cMul = binLevel(cUnary, op2("*"), op2("/"), op2("%"))
	cAdd = binLevel(cMul, op2("+"), op2("-"))
	cShift = binLevel(cAdd, op2("<<"), op2(">>"))
	cAnd = binLevel(cShift, op2("&"))
	cXor = binLevel(cAnd, op2("^"))
	cOr = binLevel(cXor, op2("|"))
	cExpr = cOr
}

// CExpr is the whole GAS expression grammar (numbers/symbols/operators).
// A function, not a var: the levels are built in init() (mutual recursion).
func CExpr() common.Combinator[rune, parsecstrings.Position, *Expr] {
	return cExpr
}

// ParseExpr parses an expression from a string (for tests and utilities);
// the whole text must be consumed.
func ParseExpr(s string) (*Expr, error) {
	body := func(buf common.Buffer[rune, parsecstrings.Position]) (*Expr, common.Error[parsecstrings.Position]) {
		e, err := cExpr(buf)
		if err != nil {
			return nil, err
		}

		if !buf.IsEOF() {
			return nil, common.NewParseError(buf.Position(), "unexpected trailing characters")
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
