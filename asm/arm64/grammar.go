package arm64

// The ARM64 instruction grammar (parsec): registers (x/w/d/h/q/s/b/v +
// sp/xzr, arrangement suffixes ".8b", the index suffix "v30[1]"),
// immediates (#expressions), floating literals, shifts/extensions,
// memory [rn, #off], lists { v0.16b, x0 }, and literal pool slots
// (=expression). It builds unevaluated operands (op.go).

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

var armRegPrefixes = "xwdbhqsv"

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

// cArmReg parses a register (letter prefix + digits + an optional
// .arr), or a named one (sp/xzr/wzr/wsp/lr/fp).
var cArmReg = func() parsec.Combinator[rune, parsecstrings.Position, armOp] {
	letter := parsecstrings.Try(parsecstrings.Satisfy("letter", true, func(r rune) bool {
		return r >= 'a' && r <= 'z'
	}))
	digit := expr.CDecDigit()
	dot := parsecstrings.Try(parsecstrings.Eq("dot", '.'))
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		var rs []rune
		first, err := letter(buf)
		if err != nil {
			return armOp{}, err
		}

		rs = append(rs, first)
		for {
			l, lerr := letter(buf)
			if lerr != nil {
				break
			}

			rs = append(rs, l)
		}

		for {
			d, derr := digit(buf)
			if derr != nil {
				break
			}

			rs = append(rs, d)
		}

		save := buf.Position()
		if _, derr := dot(buf); derr == nil {
			var arr []rune
			for {
				a, aerr := letter(buf)
				if aerr != nil {
					ad, aderr := digit(buf)
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
}()

// cArmImm parses '#' + expression or a bare expression.
var cArmImm = parsecstrings.Cast(
	parsecstrings.SkipMany(expr.CSpace,
		parsecstrings.Try(func() parsec.Combinator[rune, parsecstrings.Position, *expr.Expr] {
			return func(buf parsec.Buffer[rune, parsecstrings.Position]) (*expr.Expr, parsec.Error[parsecstrings.Position]) {
				// an optional '#' before the immediate (objdump style)
				expr.SkipHash(buf)
				return expr.CExpr()(buf)
			}
		}()),
	),
	func(e *expr.Expr) (armOp, error) {
		return newImmOp(e), nil
	},
)

// cArmFloat parses '#' + a decimal float (fmov imm).
var cArmFloat = parsecstrings.Cast(
	parsecstrings.Skip(
		parsecstrings.Try(parsecstrings.Eq("'#'", '#')),
		floatRun(),
	),
	func(rs []rune) (armOp, error) {
		v, err := strconv.ParseFloat(string(rs), 64)
		if err != nil {
			return armOp{}, err
		}

		return newFloatOp(v), nil
	},
)

// floatRun parses a decimal float: digits [. digits] ([eE]...).
func floatRun() parsec.Combinator[rune, parsecstrings.Position, []rune] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) ([]rune, parsec.Error[parsecstrings.Position]) {
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

// cArmShift parses lsl/lsr/asr/ror #imm (a modifier operand).
var cArmShift = func() parsec.Combinator[rune, parsecstrings.Position, armOp] {
	kw := parsecstrings.Try(parsecstrings.MapStrings("shift", map[string]string{
		"lsl": "lsl", "lsr": "lsr", "asr": "asr", "ror": "ror",
	}))
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		name, err := kw(buf)
		if err != nil {
			return armOp{}, err
		}

		expr.SkipSpaces(buf)
		expr.SkipHash(buf)
		var digits []rune
		digit := expr.CDecDigit()
		for {
			d, derr := digit(buf)
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
}()

// cArmExtend parses an extension (uxtw/sxtw/...) with an optional #imm.
var cArmExtend = func() parsec.Combinator[rune, parsecstrings.Position, armOp] {
	kw := parsecstrings.Try(parsecstrings.MapStrings("extend", map[string]string{
		"uxtb": "uxtb", "uxth": "uxth", "uxtw": "uxtw", "uxtx": "uxtx",
		"sxtb": "sxtb", "sxth": "sxth", "sxtw": "sxtw", "sxtx": "sxtx",
		"lsl": "lsl",
	}))
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		name, err := kw(buf)
		if err != nil {
			return armOp{}, err
		}

		op := newExtendOp(name)
		save := buf.Position()
		expr.SkipSpaces(buf)
		if _, herr := parsecstrings.Try(parsecstrings.Eq("'#'", '#'))(buf); herr == nil {
			var digits []rune
			digit := expr.CDecDigit()
			for {
				d, derr := digit(buf)
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
}()

// cArmMem parses the addressing forms: [x0] / [x0, #imm] / [x0, #imm]!
// / [x0], #imm / [x0, x1] / [x0, x1, lsl #3] / [x0, x1, uxtw].
var cArmMem = func() parsec.Combinator[rune, parsecstrings.Position, armOp] {
	lb := parsecstrings.Try(parsecstrings.Eq("'['", '['))
	rb := parsecstrings.Try(parsecstrings.Eq("']'", ']'))
	bang := parsecstrings.Try(parsecstrings.Eq("'!'", '!'))
	comma := parsecstrings.Try(parsecstrings.Eq("comma", ','))
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		m := &armMem{}
		if _, err := lb(buf); err != nil {
			return armOp{}, err
		}

		expr.SkipSpaces(buf)
		base, err := cArmReg(buf)
		if err != nil {
			return armOp{}, err
		}

		m.base = base.reg
		expr.SkipSpaces(buf)
		if _, err := rb(buf); err == nil {
			// post-index: [x0], #imm
			expr.SkipSpaces(buf)
			if _, cerr := comma(buf); cerr == nil {
				expr.SkipSpaces(buf)
				expr.SkipHash(buf)
				off, perr := expr.CExpr()(buf)
				if perr != nil {
					return armOp{}, perr
				}

				m.post = off
			}

			return newMemOp(m), nil
		}

		if _, err := comma(buf); err != nil {
			return armOp{}, parsec.NewParseError(buf.Position(), "expected ']' or ','")
		}

		expr.SkipSpaces(buf)
		// [x0, #imm]
		if r, ok := expr.PeekRune(buf); !ok || r == '#' || r == '-' || r == '0' ||
			(r >= '0' && r <= '9') || r == '(' || r == '.' || r == '\'' {
			expr.SkipHash(buf)
			off, perr := expr.CExpr()(buf)
			if perr != nil {
				return armOp{}, perr
			}

			m.offExpr = off
		} else {
			// [x0, x1(, opt #imm)]
			reg2, rerr := cArmReg(buf)
			if rerr != nil {
				return armOp{}, rerr
			}

			m.offReg = reg2.reg
			expr.SkipSpaces(buf)
			if _, oerr := comma(buf); oerr == nil {
				expr.SkipSpaces(buf)
				ext, eerr := parsecstrings.Choice("option",
					parsecstrings.Try(cArmExtend),
					parsecstrings.Try(cArmShift))(buf)
				if eerr != nil {
					return armOp{}, eerr
				}

				m.opt = ext.shift
				m.optAmt = ext.expr
			}
		}

		expr.SkipSpaces(buf)
		if _, err := rb(buf); err != nil {
			return armOp{}, err
		}

		if _, err := bang(buf); err == nil {
			m.pre = true
		}

		return newMemOp(m), nil
	}
}()

// cArmList parses { v0.16b, v1.16b } / { x0, x1 } with an optional
// index [n].
var cArmList = func() parsec.Combinator[rune, parsecstrings.Position, armOp] {
	lb := parsecstrings.Try(parsecstrings.Eq("'{'", '{'))
	rb := parsecstrings.Try(parsecstrings.Eq("'}'", '}'))
	comma := parsecstrings.Try(parsecstrings.Eq("comma", ','))
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (armOp, parsec.Error[parsecstrings.Position]) {
		if _, err := lb(buf); err != nil {
			return armOp{}, err
		}

		op := newListOp()
		for {
			expr.SkipSpaces(buf)
			r, err := cArmReg(buf)
			if err != nil {
				return armOp{}, err
			}

			op.list = append(op.list, newArmListReg(r.reg, r.arr))
			expr.SkipSpaces(buf)
			if _, cerr := comma(buf); cerr != nil {
				break
			}
		}

		expr.SkipSpaces(buf)
		if _, err := rb(buf); err != nil {
			return armOp{}, err
		}

		return op, nil
	}
}()

// cArmOperand parses an operand (the alternatives are Try-wrapped).
// cArmLitPool parses the "=expression" operand (literal pool): the
// value goes into the subsection pool, the instruction refers to the
// slot (see asm.PoolUser).
var cArmLitPool = parsecstrings.Cast(
	parsecstrings.Skip(
		parsecstrings.Try(parsecstrings.Eq("'='", '=')),
		expr.CExpr(),
	),
	func(e *expr.Expr) (armOp, error) {
		return newArmOpLit(e), nil
	},
)

var cArmOperand = parsecstrings.Choice("operand",
	parsecstrings.Try(cArmLitPool),
	parsecstrings.Try(cArmMem),
	parsecstrings.Try(cArmList),
	parsecstrings.Try(cArmFloat),
	parsecstrings.Try(cArmShift),
	parsecstrings.Try(cArmExtend),
	parsecstrings.Try(cArmReg),
	cArmImm,
)

// cArmMnemonic parses a mnemonic: [a-z][a-z0-9.]* (b.eq, ld1.16b,
// add).
var cArmMnemonic = parsecstrings.Cast(
	parsecstrings.Concat(8,
		parsecstrings.Count(1, "mnemonic",
			parsecstrings.Try(parsecstrings.Satisfy("letter", true, func(r rune) bool {
				return r >= 'a' && r <= 'z'
			}))),
		parsecstrings.Many(6,
			parsecstrings.Try(parsecstrings.Satisfy("mnemonic char", true, func(r rune) bool {
				return r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '.'
			}))),
	),
	func(rs []rune) (string, error) {
		return string(rs), nil
	},
)

// Instruction is the grammar "mnemonic comma-separated operands".
func (b *Backend) Instruction() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (asm.Unresolved, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		expr.SkipSpaces(buf)
		mnem, err := cArmMnemonic(buf)
		if err != nil {
			return nil, err
		}

		if r, ok := expr.PeekRune(buf); ok && r != ' ' && r != '\t' && r != ',' && r != '\n' {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("unknown mnemonic %q", mnem))
		}

		in := newArmAsmInstr(mnem, b.extraCtors)

		expr.SkipSpaces(buf)
		if expr.AtEOL(buf) {
			return in, nil
		}

		op, err := cArmOperand(buf)
		if err != nil {
			return nil, err
		}

		if rerr := parseRegIndex(buf, &op); rerr != nil {
			return nil, rerr
		}

		in.ops = append(in.ops, op)
		for {
			save := buf.Position()
			expr.SkipSpaces(buf)
			if _, err := expr.CComma(buf); err != nil {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				break
			}

			expr.SkipSpaces(buf)
			op, err := cArmOperand(buf)
			if err != nil {
				return nil, err
			}

			if rerr := parseRegIndex(buf, &op); rerr != nil {
				return nil, rerr
			}

			in.ops = append(in.ops, op)
		}

		return in, nil
	}
}

// parseRegIndex parses the vector register index suffix: "v8[0]". The
// suffix is optional: its absence is not an error, only a failed
// rewind is.
func parseRegIndex(
	buf parsec.Buffer[rune, parsecstrings.Position],
	op *armOp,
) parsec.Error[parsecstrings.Position] {
	if op.kind != armOpReg {
		return nil
	}

	save := buf.Position()
	lb := parsecstrings.Try(parsecstrings.Eq("'['", '['))
	rb := parsecstrings.Try(parsecstrings.Eq("']'", ']'))
	if _, err := lb(buf); err != nil {
		return expr.Rewind(buf, save)
	}

	var digits []rune
	digit := expr.CDecDigit()
	for {
		d, derr := digit(buf)
		if derr != nil {
			break
		}

		digits = append(digits, d)
	}

	if len(digits) == 0 {
		return expr.Rewind(buf, save)
	}

	if _, err := rb(buf); err != nil {
		return expr.Rewind(buf, save)
	}

	idx, err := strconv.ParseInt(string(digits), 10, 32)
	if err != nil {
		return parsec.NewParseError(buf.Position(), "register index: "+err.Error())
	}

	op.expr = expr.Num(idx)
	return nil
}

// Comment parses ';' and '//' to the end of the line ('#' is NOT a
// comment - it is the imm prefix!).
func (b *Backend) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	body := parsecstrings.Many(4, expr.CNotNL)
	semi := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.Eq("comment", ';')), body),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	slash := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.String("comment", "//")), body),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	return parsecstrings.Choice("comment", parsecstrings.Try(slash), parsecstrings.Try(semi))
}
