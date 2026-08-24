package asm

// The assembler line grammar (arch-neutral): [label:]... directive |
// instruction | empty/comment. The instruction is parsed by
// Syntax.Instruction; directives are a list of known ones with argument
// types, unknown ones are an error.

import (
	"fmt"
	"strings"

	"github.com/okneniz/parsec/common"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

// statement is one line of the source.
type statement struct {
	pos       parsecstrings.Position
	labels    []string // "label:" (there may be several)
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
	argsSecName                       // section name [.section], the rest is ignored
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

// labelC is a label definition: an identifier ':' or a numeric local
// [0-9]+ ':' (redefinable; Nb/Nf references - see resolveLocal).
var labelC = parsecstrings.Choice("label definition",
	parsecstrings.Try(labelIdentC),
	parsecstrings.Try(labelNumC),
)

// labelIdentC is a named label: identifier ':' - the name with continuation
// (digits after the first letter are legal: "foo2:").
var labelIdentC = parsecstrings.Cast(
	parsecstrings.Concat(8,
		parsecstrings.Some(4, "label name",
			parsecstrings.Try(parsecstrings.Satisfy("label start", true, expr.IsIdentStart))),
		parsecstrings.Many(8,
			parsecstrings.Try(parsecstrings.Satisfy("label char", true, expr.IsIdentCont))),
		parsecstrings.Count(1, "':'", parsecstrings.Try(parsecstrings.Eq("':'", ':'))),
	),
	func(rs []rune) (string, error) {
		return string(rs[:len(rs)-1]), nil
	},
)

// labelNumC is a numeric local label: [0-9]+ ':'.
var labelNumC = parsecstrings.Cast(
	parsecstrings.Concat(8,
		parsecstrings.Some(4, "numeric label digits", expr.CDecDigit()),
		parsecstrings.Count(1, "':'", parsecstrings.Try(parsecstrings.Eq("':'", ':'))),
	),
	func(rs []rune) (string, error) {
		return string(rs[:len(rs)-1]), nil
	},
)

// directiveC is '.' + a known directive + arguments per specification.
var directiveC = func(buf common.Buffer[rune, parsecstrings.Position]) (*directive, common.Error[parsecstrings.Position]) {
	pos := buf.Position()
	if _, err := parsecstrings.Try(parsecstrings.Eq("'.'", '.'))(buf); err != nil {
		return nil, err
	}

	name, err := cIdent(buf)
	if err != nil {
		return nil, common.NewParseError(pos, "directive name expected")
	}

	full := "." + name
	kind, ok := directives[full]
	if !ok && strings.HasPrefix(full, ".cfi_") {
		kind, ok = argsRestIgnore, true // call frame info - the whole family
	}

	if !ok {
		return nil, common.NewParseError(pos, fmt.Sprintf("unknown directive %q", full))
	}

	args, aerr := directiveArgs(kind)(buf)
	if aerr != nil {
		return nil, aerr
	}

	return newDirective(full, args), nil
}

// skipLineBody consumes everything up to end of line, NOT including the
// newline (unlike skipToEOL - parseLine will eat it).
func skipLineBody(buf common.Buffer[rune, parsecstrings.Position]) {
	for {
		if _, err := cNotNL(buf); err != nil {
			return
		}
	}
}

// directiveArgs is the directive arguments per its specification.
func directiveArgs(kind dirArgsKind) common.Combinator[rune, parsecstrings.Position, []dirArg] {
	return func(buf common.Buffer[rune, parsecstrings.Position]) ([]dirArg, common.Error[parsecstrings.Position]) {
		switch kind {
		case argsNone:
			return nil, nil
		case argsRestIgnore:
			skipLineBody(buf)
			return nil, nil
		case argsOneExpr, argsExprs:
			return exprList(buf, kind == argsOneExpr)
		case argsStrs:
			return strList(buf)
		case argsSymExpr:
			expr.SkipSpaces(buf)
			sym, err := cIdent(buf)
			if err != nil {
				return nil, err
			}

			expr.SkipSpaces(buf)
			if _, err := cComma(buf); err != nil {
				return nil, err
			}

			expr.SkipSpaces(buf)
			e, err := expr.CExpr()(buf)
			if err != nil {
				return nil, err
			}

			return []dirArg{newDirArg(nil, sym, false), newDirArg(e, "", false)}, nil
		case argsSymRest, argsSecName:
			expr.SkipSpaces(buf)
			sym, err := cIdent(buf)
			if err != nil {
				return nil, err
			}

			if kind == argsSymRest {
				skipLineBody(buf)
			}

			return []dirArg{newDirArg(nil, sym, false)}, nil
		case argsIncbin:
			// path string + optional skip, count (expressions)
			expr.SkipSpaces(buf)
			path, err := cStringLit(buf)
			if err != nil {
				return nil, err
			}

			args := []dirArg{newDirArg(nil, path, true)}
			for range 2 {
				save := buf.Position()
				expr.SkipSpaces(buf)
				if _, cerr := cComma(buf); cerr != nil {
					if rerr := expr.Rewind(buf, save); rerr != nil {
						return nil, rerr
					}

					return args, nil
				}

				expr.SkipSpaces(buf)
				e, eerr := expr.CExpr()(buf)
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
			e, err := expr.CExpr()(buf)
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
}

// exprList is comma-separated expressions; single - exactly one. Before each
// expression a '#' is allowed (objdump-style immediates: ".word #0x1234").
func exprList(
	buf common.Buffer[rune, parsecstrings.Position],
	single bool,
) ([]dirArg, common.Error[parsecstrings.Position]) {
	expr.SkipSpaces(buf)
	expr.SkipHash(buf)
	first, err := expr.CExpr()(buf)
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
		if _, err := cComma(buf); err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return nil, rerr
			}

			return args, nil
		}

		expr.SkipSpaces(buf)
		expr.SkipHash(buf)
		e, err := expr.CExpr()(buf)
		if err != nil {
			return nil, err
		}

		args = append(args, newDirArg(e, "", false))
	}
}

// strList is comma-separated string literals.
func strList(
	buf common.Buffer[rune, parsecstrings.Position],
) ([]dirArg, common.Error[parsecstrings.Position]) {
	expr.SkipSpaces(buf)
	first, err := cStringLit(buf)
	if err != nil {
		return nil, err
	}

	args := []dirArg{newDirArg(nil, first, true)}
	for {
		save := buf.Position()
		expr.SkipSpaces(buf)
		if _, err := cComma(buf); err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return nil, rerr
			}

			return args, nil
		}

		expr.SkipSpaces(buf)
		s, err := cStringLit(buf)
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
	instr := parsecstrings.Try(be.Instruction())
	comment := parsecstrings.Try(be.Comment())

	var out []statement
	for !buf.IsEOF() {
		st, err := parseLine(buf, instr, comment)
		if err != nil {
			out = append(out, newStatement(err.Position(), &[]AsmError{posErrFrom(err)}[0]))
			skipToEOL(buf)
			continue
		}

		out = append(out, st)

		// .end is the end of the source: lines below are not read at all
		// (as in GAS - not merely not assembled, their parse errors are not
		// reported either)
		if st.directive != nil && st.directive.name == ".end" {
			break
		}
	}

	return out
}

func posErrFrom(e common.Error[parsecstrings.Position]) AsmError {
	return posErr(e.Position(), e.Error())
}

// parseLine is the grammar of one line. Consumes the newline (or reaches
// EOF).
func parseLine(
	buf common.Buffer[rune, parsecstrings.Position],
	instr common.Combinator[rune, parsecstrings.Position, Unresolved],
	comment common.Combinator[rune, parsecstrings.Position, string],
) (statement, common.Error[parsecstrings.Position]) {
	start := buf.Position()
	st := newStatement(start, nil)

	expr.SkipSpaces(buf)

	// a full comment or an empty line
	if _, err := comment(buf); err == nil {
		consumeEOL(buf)
		return st, nil
	}

	if atEOL(buf) {
		consumeEOL(buf)
		return st, nil
	}

	// labels: "name:" (there may be several: "a: b: ...")
	for {
		save := buf.Position()
		expr.SkipSpaces(buf)
		lbl, err := labelC(buf)
		if err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return statement{}, rerr
			}

			break
		}

		st.labels = append(st.labels, lbl)
	}

	expr.SkipSpaces(buf)
	if atEOL(buf) {
		consumeEOL(buf)
		return st, nil
	}

	if r, ok := expr.PeekRune(buf); ok && r == '.' {
		d, err := directiveC(buf)
		if err != nil {
			return statement{}, err
		}

		st.directive = d
	} else {
		payload, err := instr(buf)
		if err != nil {
			return statement{}, err
		}

		st.instr = payload
		st.hasInstr = true
	}

	expr.SkipSpaces(buf)
	consumeComment(buf, comment)
	if !atEOL(buf) {
		return statement{}, common.NewParseError(buf.Position(), "unexpected trailing characters")
	}

	consumeEOL(buf)
	return st, nil
}
