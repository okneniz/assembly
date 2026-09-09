package asm

// Lexical combinators of the assembler line (parsec/strings): newline,
// comma, identifiers (directives/labels), string literals. The base rune
// atoms (spaces/digits/name predicates/peek/rewind) and all expression
// lexing are in asm/expr: the expression is the lowest text layer, the line
// grammar is built on it. All atoms are Try-wrapped: a failed greedy atom
// leaves the position advanced, restoring it is the caller alternative's
// job (the parsec library contract).
//
// The combinators are values built by their constructors (no package
// vars); the line grammar (stmt.go) captures them once per source.

import (
	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

// newIdent is an identifier/directive or label name:
// [._$a-zA-Z][._$a-zA-Z0-9]* - the continuation may contain digits
// ("p2align", "foo2").
func makeIdentParser() parsec.Combinator[rune, parsecstrings.Position, string] {
	return parsecstrings.Cast(
		parsecstrings.Concat(8,
			parsecstrings.Some(
				4,
				"identifier start",
				parsecstrings.Try(
					parsecstrings.Satisfy("identifier start", true, expr.IsIdentStart),
				),
			),
			parsecstrings.Many(
				8,
				parsecstrings.Try(
					parsecstrings.Satisfy("identifier char", true, expr.IsIdentCont),
				),
			),
		),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
}

// newStringLit is a "..." string literal with escape sequences.
func makeStringLitParser() parsec.Combinator[rune, parsecstrings.Position, string] {
	dquote := parsecstrings.Try(parsecstrings.Eq("double quote", '"'))

	return parsecstrings.Cast(
		parsecstrings.Between(dquote, makeStringBody(), dquote),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
}

// stringBody is any characters except '"' and '\n'; escape sequences are
// expanded.
func makeStringBody() parsec.Combinator[rune, parsecstrings.Position, []rune] {
	anyRune := parsecstrings.Any()

	return func(buf parsec.Buffer[rune, parsecstrings.Position]) ([]rune, parsec.Error[parsecstrings.Position]) {
		var out []rune
		for {
			pos := buf.Position()
			r, err := anyRune(buf)
			if err != nil {
				return nil, parsec.NewParseError(pos, "unterminated string")
			}

			if r == '"' {
				if err := buf.Seek(pos); err != nil {
					return nil, parsec.NewParseError(pos, err.Error())
				}

				return out, nil
			}

			if r == '\n' {
				return nil, parsec.NewParseError(pos, "newline in string literal")
			}

			if r == '\\' {
				e, err := anyRune(buf)
				if err != nil {
					return nil, parsec.NewParseError(pos, "unterminated escape")
				}

				switch e {
				case 'n':
					r = '\n'
				case 't':
					r = '\t'
				case 'r':
					r = '\r'
				case '0':
					r = 0
				default: // \\ \" and the rest - as is
					r = e
				}
			}

			out = append(out, r)
		}
	}
}

// atEOL is true when the end of line follows (newline or EOF).
func atEOL(buf parsec.Buffer[rune, parsecstrings.Position]) bool {
	return expr.AtEOL(buf)
}

// consumeEOL consumes '\n' if present.
func consumeEOL(buf parsec.Buffer[rune, parsecstrings.Position]) {
	if r, ok := expr.PeekRune(buf); ok && r == '\n' {
		if err := expr.ConsumeRune(buf); err != nil {
			return // the rune just peeked - unreadable only at I/O failure
		}
	}
}

// consumeComment consumes a trailing comment if present (the combinator is
// passed by the backend); the absence of a comment is not an error.
func consumeComment(
	buf parsec.Buffer[rune, parsecstrings.Position],
	c parsec.Combinator[rune, parsecstrings.Position, string],
) {
	if _, err := c(buf); err != nil {
		return // no comment - nothing to consume
	}
}

// skipToEOL consumes everything up to end of line inclusive (recovery after
// a line parse error).
func skipToEOL(buf parsec.Buffer[rune, parsecstrings.Position]) {
	skipLineBody(buf)
	consumeEOL(buf)
}

// skipLineBody consumes everything up to end of line, NOT including the
// newline (unlike skipToEOL - parseLine will eat it).
func skipLineBody(buf parsec.Buffer[rune, parsecstrings.Position]) {
	for {
		r, ok := expr.PeekRune(buf)
		if !ok || r == '\n' {
			return
		}

		if err := expr.ConsumeRune(buf); err != nil {
			return // the rune just peeked - unreadable only at I/O failure
		}
	}
}
