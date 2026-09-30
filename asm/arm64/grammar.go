package arm64

// The ARM64 instruction grammar: a mnemonic, then comma-separated
// operands - registers (with .arr and [idx]), immediates/expressions,
// memory addressing forms, register lists, literals, shifts and
// extends. It builds unevaluated armOp values; the instruction itself
// (armAsmInstr) is assembled by the ctors (asm_ctors.go) and encoded by
// encodeARM (encode.go).
//
// The grammar is a struct built once per Backend (makeGrammar): the
// combinators are values captured by the closures; the memory-addressing
// branch refers to the register/extend/shift operands through the struct
// fields (lazily - all fields are assigned before the first parse).

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

const armRegPrefixes = "xwdbhqsv"

// isRegisterName checks whether s looks like a register name (x0, w31,
// sp, v3.8b...).
func isRegisterName(s string) bool {
	switch s {
	case "sp", "xzr", "wzr", "wsp", "lr", "fp":
		return true
	}

	if len(s) < 2 || !strings.ContainsRune(armRegPrefixes, rune(s[0])) {
		return false
	}

	num := s[1:]
	if i := strings.IndexByte(num, '.'); i >= 0 {
		num = num[:i]
	}

	if num == "" {
		return false
	}

	for _, r := range num {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

type armOperand = parsec.Combinator[rune, parsecstrings.Position, armOp, parsec.Stateless]

// armGrammar is the whole operand grammar plus the shared atoms (the
// expression ladder, the decimal digit, the punctuation) captured once.
type armGrammar struct {
	parseReg      armOperand
	parseImm      armOperand
	parseFloat    armOperand
	parseShift    armOperand
	parseExtend   armOperand
	parseMem      armOperand
	parseList     armOperand
	parseLitPool  armOperand
	parseOperand  armOperand
	parseMnemonic parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless]

	parseExpr   parsec.Combinator[rune, parsecstrings.Position, *expr.Expr, parsec.Stateless]
	parseDigit  parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseComma  parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseHash   parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseLBrack parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseRBrack parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseLBrace parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseRBrace parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
}

// makeGrammar builds the whole grammar once; the ready combinators
// are captured by the backend and reused for every line.
func makeGrammar() *armGrammar {
	g := &armGrammar{
		parseExpr:   expr.MakeExprParser(),
		parseDigit:  expr.MakeDigitParser(),
		parseComma:  expr.MakeCommaParser(),
		parseHash:   parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'#'", '#')),
		parseLBrack: parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'['", '[')),
		parseRBrack: parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("']'", ']')),
		parseLBrace: parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'{'", '{')),
		parseRBrace: parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'}'", '}')),
	}

	letter := parsecstrings.Try(parsecstrings.Satisfy[parsec.Stateless]("letter", true, func(r rune) bool {
		return r >= 'a' && r <= 'z'
	}))
	dot := parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("dot", '.'))

	// a register (letter prefix + digits + an optional .arr), or a named
	// one (sp/xzr/wzr/wsp/lr/fp)
	g.parseReg = func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		var rs []rune
		first, err := letter(state, buf)
		if err != nil {
			return armOp{}, err
		}

		rs = append(rs, first)
		for {
			l, lerr := letter(state, buf)
			if lerr != nil {
				break
			}

			rs = append(rs, l)
		}

		for {
			d, derr := g.parseDigit(state, buf)
			if derr != nil {
				break
			}

			rs = append(rs, d)
		}

		save := buf.Position()
		if _, derr := dot(state, buf); derr == nil {
			var arr []rune
			for {
				a, aerr := letter(state, buf)
				if aerr != nil {
					ad, aderr := g.parseDigit(state, buf)
					if aderr != nil {
						break
					}

					arr = append(arr, ad)
					continue
				}

				arr = append(arr, a)
			}

			if len(arr) == 0 {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return armOp{}, rerr
				}
			} else {
				rs = append(rs, '.')
				rs = append(rs, arr...)
			}
		}

		name := string(rs)
		// an identifier continues past the token (sp_el0, cntvct_el0):
		// a register never continues with '_' - give the token back to
		// the expression lexer (Try rewinds the position)
		if r, ok := expr.PeekRune(buf); ok && r == '_' {
			return armOp{}, parsec.NewParseError(
				buf.Position(),
				fmt.Sprintf("not a register: %q", name),
			)
		}

		if !isRegisterName(name) {
			// the position did not move - no rewind needed, the error
			// goes to the caller's Try
			return armOp{}, parsec.NewParseError(
				buf.Position(),
				fmt.Sprintf("not a register: %q", name),
			)
		}

		op := newRegOp(name)
		if before, after, ok := strings.Cut(name, "."); ok {
			op.reg, op.arr = before, after
		}

		return op, nil
	}

	// '#' + expression or a bare expression
	g.parseImm = parsecstrings.Cast(
		parsecstrings.SkipMany(
			expr.MakeSpaceParser(),
			parsecstrings.Try(
				func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (*expr.Expr, parsec.Error[parsecstrings.Position]) {
					// an optional '#' before the immediate (objdump style)
					expr.SkipHash(buf)
					return g.parseExpr(state, buf)
				},
			),
		),
		func(e *expr.Expr) (armOp, error) {
			return newImmOp(e), nil
		},
	)

	// '#' + a decimal float (fmov imm)
	g.parseFloat = parsecstrings.Cast(
		parsecstrings.Skip(g.parseHash, makeFloatRun()),
		func(rs []rune) (armOp, error) {
			v, err := strconv.ParseFloat(string(rs), 64)
			if err != nil {
				return armOp{}, err
			}

			return newFloatOp(v), nil
		},
	)

	// lsl/lsr/asr/ror #imm (a modifier operand)
	shiftKw := parsecstrings.Try(parsecstrings.MapStrings[string, parsec.Stateless]("shift", map[string]string{
		"lsl": "lsl", "lsr": "lsr", "asr": "asr", "ror": "ror",
	}))
	g.parseShift = func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		name, err := shiftKw(state, buf)
		if err != nil {
			return armOp{}, err
		}

		expr.SkipSpaces(buf)
		expr.SkipHash(buf)
		var digits []rune
		for {
			d, derr := g.parseDigit(state, buf)
			if derr != nil {
				break
			}

			digits = append(digits, d)
		}

		if len(digits) == 0 {
			return armOp{}, parsec.NewParseError(buf.Position(), "shift amount expected")
		}

		amt, aerr := strconv.ParseInt(string(digits), 10, 32)
		if aerr != nil {
			return armOp{}, parsec.NewParseError(buf.Position(), "shift amount: "+aerr.Error())
		}

		return newShiftOp(name, expr.Num(amt)), nil
	}

	// an extension (uxtw/sxtw/...) with an optional #imm
	extendKw := parsecstrings.Try(parsecstrings.MapStrings[string, parsec.Stateless]("extend", map[string]string{
		"uxtb": "uxtb", "uxth": "uxth", "uxtw": "uxtw", "uxtx": "uxtx",
		"sxtb": "sxtb", "sxth": "sxth", "sxtw": "sxtw", "sxtx": "sxtx",
		"lsl": "lsl",
	}))
	g.parseExtend = func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		name, err := extendKw(state, buf)
		if err != nil {
			return armOp{}, err
		}

		op := newExtendOp(name)
		save := buf.Position()
		expr.SkipSpaces(buf)
		if _, herr := g.parseHash(state, buf); herr == nil {
			var digits []rune
			for {
				d, derr := g.parseDigit(state, buf)
				if derr != nil {
					break
				}

				digits = append(digits, d)
			}

			if len(digits) > 0 {
				amt, aerr := strconv.ParseInt(string(digits), 10, 32)
				if aerr != nil {
					return armOp{}, parsec.NewParseError(
						buf.Position(),
						"extend amount: "+aerr.Error(),
					)
				}

				op.expr = expr.Num(amt)
			}
		} else {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return armOp{}, rerr
			}
		}

		return op, nil
	}

	bang := parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'!'", '!'))
	memOption := parsecstrings.Choice(
		"option",
		parsecstrings.Try(
			func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
				return g.parseExtend(state, buf)
			},
		),
		parsecstrings.Try(
			func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
				return g.parseShift(state, buf)
			},
		),
	)

	// the addressing forms: [x0] / [x0, #imm] / [x0, #imm]! / [x0], #imm
	// / [x0, x1] / [x0, x1, lsl #3] / [x0, x1, uxtw]
	g.parseMem = func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		m := &armMem{}
		if _, err := g.parseLBrack(state, buf); err != nil {
			return armOp{}, err
		}

		expr.SkipSpaces(buf)
		base, err := g.parseReg(state, buf)
		if err != nil {
			return armOp{}, err
		}

		m.base = base.reg
		expr.SkipSpaces(buf)
		if _, err := g.parseRBrack(state, buf); err == nil {
			// post-index: [x0], #imm
			expr.SkipSpaces(buf)
			if _, cerr := g.parseComma(state, buf); cerr == nil {
				expr.SkipSpaces(buf)
				expr.SkipHash(buf)
				off, perr := g.parseExpr(state, buf)
				if perr != nil {
					return armOp{}, perr
				}

				m.post = off
			}

			return newMemOp(m), nil
		}

		if _, err := g.parseComma(state, buf); err != nil {
			return armOp{}, parsec.NewParseError(buf.Position(), "expected ']' or ','")
		}

		expr.SkipSpaces(buf)
		// [x0, #imm]
		if r, ok := expr.PeekRune(buf); !ok || r == '#' || r == '-' ||
			(r >= '0' && r <= '9') || r == '(' || r == '.' || r == '\'' {
			expr.SkipHash(buf)
			off, perr := g.parseExpr(state, buf)
			if perr != nil {
				return armOp{}, perr
			}

			m.offExpr = off
		} else {
			// [x0, x1(, opt #imm)]
			reg2, rerr := g.parseReg(state, buf)
			if rerr != nil {
				return armOp{}, rerr
			}

			m.offReg = reg2.reg
			expr.SkipSpaces(buf)
			if _, oerr := g.parseComma(state, buf); oerr == nil {
				expr.SkipSpaces(buf)
				ext, eerr := memOption(state, buf)
				if eerr != nil {
					return armOp{}, eerr
				}

				m.opt = ext.shift
				m.optAmt = ext.expr
			}
		}

		expr.SkipSpaces(buf)
		if _, err := g.parseRBrack(state, buf); err != nil {
			return armOp{}, err
		}

		if _, err := bang(state, buf); err == nil {
			m.pre = true
		}

		return newMemOp(m), nil
	}

	// { v0.16b, v1.16b } / { x0, x1 }
	g.parseList = func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		if _, err := g.parseLBrace(state, buf); err != nil {
			return armOp{}, err
		}

		op := newListOp()
		for {
			expr.SkipSpaces(buf)
			r, err := g.parseReg(state, buf)
			if err != nil {
				return armOp{}, err
			}

			op.list = append(op.list, newArmListReg(r.reg, r.arr))
			expr.SkipSpaces(buf)
			if _, cerr := g.parseComma(state, buf); cerr != nil {
				break
			}
		}

		expr.SkipSpaces(buf)
		if _, err := g.parseRBrace(state, buf); err != nil {
			return armOp{}, err
		}

		return op, nil
	}

	// "=expression" (literal pool): the value goes into the subsection
	// pool, the instruction refers to the slot (see asm.PoolUser)
	g.parseLitPool = parsecstrings.Cast(
		parsecstrings.Skip(
			parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'='", '=')),
			g.parseExpr,
		),
		func(e *expr.Expr) (armOp, error) {
			return newArmOpLit(e), nil
		},
	)

	g.parseOperand = parsecstrings.Choice("operand",
		parsecstrings.Try(g.parseLitPool),
		parsecstrings.Try(g.parseMem),
		parsecstrings.Try(g.parseList),
		parsecstrings.Try(g.parseFloat),
		parsecstrings.Try(g.parseShift),
		parsecstrings.Try(g.parseExtend),
		parsecstrings.Try(g.parseReg),
		g.parseImm,
	)

	// a parseMnemonic: [a-z][a-z0-9.]* (b.eq, ld1.16b, add)
	g.parseMnemonic = parsecstrings.Cast(
		parsecstrings.Concat(8,
			parsecstrings.Count(1, "mnemonic",
				parsecstrings.Try(parsecstrings.Satisfy[parsec.Stateless]("letter", true, func(r rune) bool {
					return r >= 'a' && r <= 'z'
				}))),
			parsecstrings.Many(6,
				parsecstrings.Try(parsecstrings.Satisfy[parsec.Stateless]("mnemonic char", true, func(r rune) bool {
					return r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '.'
				}))),
		),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)

	return g
}

// floatRun parses a decimal parseFloat: digits [. digits] ([eE]...).
func makeFloatRun() parsec.Combinator[rune, parsecstrings.Position, []rune, parsec.Stateless] {
	return func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) ([]rune, parsec.Error[parsecstrings.Position]) {
		var out []rune
		digits := func() bool {
			n := 0
			for {
				r, ok := expr.PeekRune(buf)
				if !ok || r < '0' || r > '9' {
					break
				}

				out = append(out, r)
				if rerr := expr.ConsumeRune(buf); rerr != nil {
					return false
				}

				n++
			}

			return n > 0
		}
		if !digits() {
			return nil, parsec.NewParseError(buf.Position(), "float expected")
		}

		if r, ok := expr.PeekRune(buf); ok && r == '.' {
			out = append(out, r)
			if rerr := expr.ConsumeRune(buf); rerr != nil {
				return nil, rerr
			}

			digits()
		}

		if r, ok := expr.PeekRune(buf); ok && (r == 'e' || r == 'E') {
			save := buf.Position()
			out = append(out, r)
			if rerr := expr.ConsumeRune(buf); rerr != nil {
				return nil, rerr
			}

			if r2, ok2 := expr.PeekRune(buf); ok2 && (r2 == '+' || r2 == '-') {
				out = append(out, r2)
				if rerr := expr.ConsumeRune(buf); rerr != nil {
					return nil, rerr
				}
			}

			if !digits() {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				out = out[:len(out)-1]
			}
		}

		if !strings.ContainsAny(string(out), ".eE") {
			return nil, parsec.NewParseError(buf.Position(), "float expected")
		}

		return out, nil
	}
}
