package asm

// The assembler line grammar (arch-neutral): [parseLabel:]... directive |
// instruction | empty/comment. The instruction is parsed by
// Syntax.Instruction; directives are a list of known ones with argument
// types, unknown ones are an error.

import (
	"fmt"
	"strings"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

// statement is one line of the source.
type statement struct {
	pos       parsecstrings.Position
	labels    []string // "parseLabel:" (there may be several)
	directive *directive
	instr     Unresolved // unresolved instruction from the Syntax grammar
	hasInstr  bool       // there was an instruction (instr may be a typed nil)
	err       *AsmError  // line parse error (the line is skipped)
}

func newStatement(pos parsecstrings.Position, err *AsmError) statement {
	return statement{
		pos: pos,
		err: err,
	}
}

// directive is a recognized directive with arguments.
type directive struct {
	name string
	args []dirArg
}

func newDirective(name string, args []dirArg) *directive {
	return &directive{
		name: name,
		args: args,
	}
}

// dirArg is a directive argument: an expression or a string literal.
type dirArg struct {
	expr  *expr.Expr
	str   string
	isStr bool
}

func newDirArg(e *expr.Expr, str string, isStr bool) dirArg {
	return dirArg{
		expr:  e,
		str:   str,
		isStr: isStr,
	}
}

// dirArgsKind is which arguments a directive accepts.
type dirArgsKind int

const (
	argsNone       dirArgsKind = iota // no arguments
	argsExprs                         // comma-separated expression list
	argsOneExpr                       // exactly one expression
	argsStrs                          // comma-separated string list
	argsSymExpr                       // symbol, comma, expression (.set)
	argsSymRest                       // symbol, the rest of the line is ignored (.type/.size)
	argsSecName                       // section name, then ignored quoted flags [.section]
	argsRestIgnore                    // the whole argument line is ignored (.file, .cfi_*, ...)
	argsIncbin                        // path string + optional skip/count (.incbin)
	argsSubsec                        // optional subsection number (.text/.data/.bss)
)

var directives = map[string]dirArgsKind{
	".text": argsSubsec, ".data": argsSubsec, ".bss": argsSubsec,
	".word": argsExprs, ".half": argsExprs, ".short": argsExprs,
	".byte": argsExprs, ".quad": argsExprs, ".dword": argsExprs,
	".zero": argsOneExpr, ".space": argsOneExpr, ".skip": argsOneExpr,
	".end":   argsNone,
	".ltorg": argsNone,
	".align": argsOneExpr, ".p2align": argsOneExpr, ".balign": argsOneExpr,
	".string": argsStrs, ".asciz": argsStrs, ".ascii": argsStrs,
	".incbin": argsIncbin,
	".global": argsRestIgnore, ".globl": argsRestIgnore, ".local": argsRestIgnore,
	".set": argsSymExpr, ".equ": argsSymExpr,
	".type": argsSymRest, ".size": argsSymRest,
	".section": argsSecName,
	// debug/metadata - recognized and silently ignored (compiler output);
	// .option is semantic (encoding modes), applied by Syntax
	".file": argsRestIgnore, ".loc": argsRestIgnore, ".ident": argsRestIgnore,
	".attribute": argsRestIgnore, ".abicalls": argsRestIgnore, ".nan": argsRestIgnore,
	".module": argsRestIgnore, ".p2align_manual": argsRestIgnore,
	".option": argsSymRest,
}

// lineGrammar is the line-level grammar state built once per source
// (parseSource): the instruction/comment combinators of the Syntax
// backend (Try-wrapped), the label and directive grammars, the shared
// atoms (identifier, string literal, comma) and the expression ladder
// (expr.CExpr is a fresh ladder per call - capturing it once is the
// point). A struct: everything below captures its fields. sep is the
// backend's statement separator (0: one statement per line).
type lineGrammar struct {
	parseInstruction parsec.Combinator[rune, parsecstrings.Position, Unresolved, parsec.Stateless]
	parseComment     parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless]
	parseLabel       parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless]
	parseDirective   parsec.Combinator[rune, parsecstrings.Position, *directive, parsec.Stateless]
	parseIdent       parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless]
	parseStringLit   parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless]
	parseComma       parsec.Combinator[rune, parsecstrings.Position, rune, parsec.Stateless]
	parseExpr        parsec.Combinator[rune, parsecstrings.Position, *expr.Expr, parsec.Stateless]
	sep              rune
}

// makeLineGrammar builds the whole line grammar once; be is the syntax
// backend (its Instruction/Comment are captured here, not rebuilt).
func makeLineGrammar(be Syntax) *lineGrammar {
	g := &lineGrammar{
		parseInstruction: parsecstrings.Try(be.Instruction()),
		parseComment:     parsecstrings.Try(be.Comment()),
		parseIdent:       makeIdentParser(),
		parseStringLit:   makeStringLitParser(),
		parseComma:       expr.MakeCommaParser(),
		parseExpr:        expr.MakeExprParser(),
		sep:              be.Separator(),
	}

	identColon := func() parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless] {
		colon := parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("':'", ':'))
		return parsecstrings.Cast(
			parsecstrings.Concat(8,
				parsecstrings.Some(
					4,
					"label name",
					parsecstrings.Try(
						parsecstrings.Satisfy[parsec.Stateless]("label start", true, expr.IsIdentStart),
					),
				),
				parsecstrings.Many(8,
					parsecstrings.Try(parsecstrings.Satisfy[parsec.Stateless]("label char", true, expr.IsIdentCont))),
				parsecstrings.Count(1, "':'", colon),
			),
			func(rs []rune) (string, error) {
				return string(rs[:len(rs)-1]), nil
			},
		)
	}

	labelIdent := identColon()

	numColon := parsecstrings.Cast(
		parsecstrings.Concat(8,
			parsecstrings.Some(4, "numeric label digits", expr.MakeDigitParser()),
			parsecstrings.Count(1, "':'", parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("':'", ':'))),
		),
		func(rs []rune) (string, error) {
			return string(rs[:len(rs)-1]), nil
		},
	)

	// label definition: an identifier ':' or a numeric local [0-9]+ ':'
	// (redefinable; Nb/Nf references - see resolveLocal)
	g.parseLabel = parsecstrings.Choice("label definition",
		parsecstrings.Try(labelIdent),
		parsecstrings.Try(numColon),
	)

	g.parseDirective = g.makeDirectiveParser()

	return g
}

// newDirective is '.' + a known directive + arguments per specification.
func (g *lineGrammar) makeDirectiveParser() parsec.Combinator[rune, parsecstrings.Position, *directive, parsec.Stateless] {
	dot := parsecstrings.Try(parsecstrings.Eq[parsec.Stateless]("'.'", '.'))

	return func(state parsec.Stateless, buf parsec.Buffer[rune, parsecstrings.Position]) (*directive, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		if _, err := dot(state, buf); err != nil {
			return nil, err
		}

		name, err := g.parseIdent(state, buf)
		if err != nil {
			return nil, parsec.NewParseError(pos, "directive name expected")
		}

		full := "." + name
		kind, ok := directives[full]
		if !ok && strings.HasPrefix(full, ".cfi_") {
			kind, ok = argsRestIgnore, true // call frame info - the whole family
		}

		if !ok {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("unknown directive %q", full))
		}

		args, aerr := g.parseArgs(buf, kind)
		if aerr != nil {
			return nil, aerr
		}

		return newDirective(full, args), nil
	}
}

// parseArgs is the directive arguments per its specification.
func (g *lineGrammar) parseArgs(
	buf parsec.Buffer[rune, parsecstrings.Position],
	kind dirArgsKind,
) ([]dirArg, parsec.Error[parsecstrings.Position]) {
	switch kind {
	case argsNone:
		return nil, nil
	case argsRestIgnore:
		g.skipStmtBody(buf)
		return nil, nil
	case argsOneExpr, argsExprs:
		return g.exprList(buf, kind == argsOneExpr)
	case argsStrs:
		return g.strList(buf)
	case argsSymExpr:
		expr.SkipSpaces(buf)
		sym, err := g.parseIdent(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		expr.SkipSpaces(buf)
		if _, err := g.parseComma(parsec.Stateless{}, buf); err != nil {
			return nil, err
		}

		expr.SkipSpaces(buf)
		e, err := g.parseExpr(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		return []dirArg{newDirArg(nil, sym, false), newDirArg(e, "", false)}, nil
	case argsSymRest:
		expr.SkipSpaces(buf)
		sym, err := g.parseIdent(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		g.skipStmtBody(buf)
		return []dirArg{newDirArg(nil, sym, false)}, nil
	case argsSecName:
		// the section name, then the GAS flags - a comma list of quoted
		// strings (or a numeric subsection) - recognized and ignored: the
		// core carries no flag semantics
		expr.SkipSpaces(buf)
		sym, err := g.parseIdent(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		for {
			save := buf.Position()
			expr.SkipSpaces(buf)
			if _, cerr := g.parseComma(parsec.Stateless{}, buf); cerr != nil {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				break
			}

			expr.SkipSpaces(buf)
			if _, serr := g.parseStringLit(parsec.Stateless{}, buf); serr != nil {
				return nil, serr
			}
		}

		return []dirArg{newDirArg(nil, sym, false)}, nil
	case argsIncbin:
		// path string + optional skip, count (expressions)
		expr.SkipSpaces(buf)
		path, err := g.parseStringLit(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		args := []dirArg{newDirArg(nil, path, true)}
		for range 2 {
			save := buf.Position()
			expr.SkipSpaces(buf)
			if _, cerr := g.parseComma(parsec.Stateless{}, buf); cerr != nil {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				return args, nil
			}

			expr.SkipSpaces(buf)
			e, eerr := g.parseExpr(parsec.Stateless{}, buf)
			if eerr != nil {
				return nil, eerr
			}

			args = append(args, newDirArg(e, "", false))
		}

		return args, nil
	case argsSubsec:
		// optional subsection number (0..8192, as in GAS)
		save := buf.Position()
		expr.SkipSpaces(buf)
		e, err := g.parseExpr(parsec.Stateless{}, buf)
		if err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return nil, rerr
			}

			return nil, nil // ".text" without an argument - subsection 0
		}

		return []dirArg{newDirArg(e, "", false)}, nil
	}

	return nil, nil
}

// exprList is comma-separated expressions; single - exactly one. Before each
// expression a '#' is allowed (objdump-style immediates: ".word #0x1234").
func (g *lineGrammar) exprList(
	buf parsec.Buffer[rune, parsecstrings.Position],
	single bool,
) ([]dirArg, parsec.Error[parsecstrings.Position]) {
	expr.SkipSpaces(buf)
	expr.SkipHash(buf)
	first, err := g.parseExpr(parsec.Stateless{}, buf)
	if err != nil {
		return nil, err
	}

	args := []dirArg{newDirArg(first, "", false)}
	if single {
		return args, nil
	}

	for {
		save := buf.Position()
		expr.SkipSpaces(buf)
		if _, err := g.parseComma(parsec.Stateless{}, buf); err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return nil, rerr
			}

			return args, nil
		}

		expr.SkipSpaces(buf)
		expr.SkipHash(buf)
		e, err := g.parseExpr(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		args = append(args, newDirArg(e, "", false))
	}
}

// strList is comma-separated string literals.
func (g *lineGrammar) strList(
	buf parsec.Buffer[rune, parsecstrings.Position],
) ([]dirArg, parsec.Error[parsecstrings.Position]) {
	expr.SkipSpaces(buf)
	first, err := g.parseStringLit(parsec.Stateless{}, buf)
	if err != nil {
		return nil, err
	}

	args := []dirArg{newDirArg(nil, first, true)}
	for {
		save := buf.Position()
		expr.SkipSpaces(buf)
		if _, err := g.parseComma(parsec.Stateless{}, buf); err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return nil, rerr
			}

			return args, nil
		}

		expr.SkipSpaces(buf)
		s, err := g.parseStringLit(parsec.Stateless{}, buf)
		if err != nil {
			return nil, err
		}

		args = append(args, newDirArg(nil, s, true))
	}
}

// parseSource splits the source into statements; lines with parse errors get
// err and are skipped (all errors are returned in the slice). The Syntax
// grammars are wrapped in Try: the parsec contract is that a failed greedy
// combinator may leave the position advanced, restoration is explicit here.
func parseSource(src []rune, be Syntax) []statement {
	buf := parsecstrings.Buffer(src)
	g := makeLineGrammar(be)

	var out []statement
	for !buf.IsEOF() {
		sts, err := parseLine(buf, g)
		if err != nil {
			out = append(out, newStatement(err.Position(), newPosErr(err)))
			skipToEOL(buf)
			continue
		}

		out = append(out, sts...)

		// .end is the end of the source: lines below are not read at all
		// (as in GAS - not merely not assembled, their parse errors are not
		// reported either)
		for _, st := range sts {
			if st.directive != nil && st.directive.name == ".end" {
				return out
			}
		}
	}

	return out
}

func posErrFrom(e parsec.Error[parsecstrings.Position]) AsmError {
	return posErr(e.Position(), e.Error())
}

// newPosErr is posErrFrom returning a pointer (the statement error slot).
func newPosErr(e parsec.Error[parsecstrings.Position]) *AsmError {
	err := posErrFrom(e)
	return &err
}

// atSep — the backend's statement separator follows (0: the backend has
// none).
func (g *lineGrammar) atSep(buf parsec.Buffer[rune, parsecstrings.Position]) bool {
	if g.sep == 0 {
		return false
	}

	r, ok := expr.PeekRune(buf)
	return ok && r == g.sep
}

// parseLine is the grammar of one line: one or more statements separated
// by the backend's separator (a single statement when it is 0). Consumes
// the newline (or reaches EOF); a parse error kills the rest of the line.
func parseLine(
	buf parsec.Buffer[rune, parsecstrings.Position],
	g *lineGrammar,
) ([]statement, parsec.Error[parsecstrings.Position]) {
	var out []statement
	for {
		st := newStatement(buf.Position(), nil)

		expr.SkipSpaces(buf)

		// a full comment or an empty line
		if _, err := g.parseComment(parsec.Stateless{}, buf); err == nil {
			consumeEOL(buf)
			return append(out, st), nil
		}

		if atEOL(buf) {
			consumeEOL(buf)
			return append(out, st), nil
		}

		// labels: "name:" (there may be several: "a: b: ...")
		for {
			save := buf.Position()
			expr.SkipSpaces(buf)
			lbl, err := g.parseLabel(parsec.Stateless{}, buf)
			if err != nil {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				break
			}

			st.labels = append(st.labels, lbl)
		}

		expr.SkipSpaces(buf)
		if atEOL(buf) {
			consumeEOL(buf)
			return append(out, st), nil
		}

		// a label-only (or empty) segment before the separator
		if g.atSep(buf) {
			_ = expr.ConsumeRune(buf)
			out = append(out, st)
			continue
		}

		if r, ok := expr.PeekRune(buf); ok && r == '.' {
			d, err := g.parseDirective(parsec.Stateless{}, buf)
			if err != nil {
				return nil, err
			}

			st.directive = d
		} else {
			payload, err := g.parseInstruction(parsec.Stateless{}, buf)
			if err != nil {
				return nil, err
			}

			st.instr = payload
			st.hasInstr = true
		}

		expr.SkipSpaces(buf)
		consumeComment(buf, g.parseComment)
		if atEOL(buf) {
			consumeEOL(buf)
			return append(out, st), nil
		}

		if g.atSep(buf) {
			_ = expr.ConsumeRune(buf)
			out = append(out, st)
			continue
		}

		return nil, parsec.NewParseError(buf.Position(), "unexpected trailing characters")
	}
}

// skipStmtBody consumes everything up to the statement separator or end
// of line, quote-aware (a separator inside a string literal is data),
// NOT consuming the newline/separator itself - parseLine eats them. The
// ignored directive tails stop before a next statement begins.
func (g *lineGrammar) skipStmtBody(buf parsec.Buffer[rune, parsecstrings.Position]) {
	inStr := false
	for {
		r, ok := expr.PeekRune(buf)
		if !ok || r == '\n' {
			return
		}

		if !inStr && g.sep != 0 && r == g.sep {
			return
		}

		if r == '"' {
			inStr = !inStr
		}

		if err := expr.ConsumeRune(buf); err != nil {
			return // the rune just peeked - unreadable only at I/O failure
		}
	}
}
