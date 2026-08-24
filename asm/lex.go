package asm

// Lexical combinators of the assembler line (parsec/strings): newline,
// comma, identifiers (directives/labels), string literals. The base rune
// atoms (spaces/digits/name predicates/peek/rewind) and all expression
// lexing are in asm/expr: the expression is the lowest text layer, the line
// grammar is built on it. All atoms are Try-wrapped: a failed greedy atom
// leaves the position advanced, restoring it is the caller alternative's
// job (the parsec library contract).

import (
	"github.com/okneniz/parsec/common"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

var (
	cNotNL  = expr.CNotNL
	cNL     = expr.CNL
	cComma  = expr.CComma
	cDQuote = parsecstrings.Try(parsecstrings.Eq("double quote", '"'))
)

// cIdent is an identifier/directive or label name:
// [._$a-zA-Z][._$a-zA-Z0-9]* - the continuation may contain digits
// ("p2align", "foo2").
var cIdent = parsecstrings.Cast(
	parsecstrings.Concat(8,
		parsecstrings.Some(4, "identifier start",
			parsecstrings.Try(parsecstrings.Satisfy("identifier start", true, expr.IsIdentStart))),
		parsecstrings.Many(8,
			parsecstrings.Try(parsecstrings.Satisfy("identifier char", true, expr.IsIdentCont))),
	),
	func(rs []rune) (string, error) {
		return string(rs), nil
	},
)

// cStringLit is a "..." string literal with escape sequences.
var cStringLit = parsecstrings.Cast(
	parsecstrings.Between(cDQuote, stringBody(), cDQuote),
	func(rs []rune) (string, error) {
		return string(rs), nil
	},
)

// stringBody is any characters except '"' and '\n'; escape sequences are
// expanded.
func stringBody() common.Combinator[rune, parsecstrings.Position, []rune] {
	return func(buf common.Buffer[rune, parsecstrings.Position]) ([]rune, common.Error[parsecstrings.Position]) {
		var out []rune
		for {
			pos := buf.Position()
			r, err := parsecstrings.Any()(buf)
			if err != nil {
				return nil, common.NewParseError(pos, "unterminated string")
			}

			if r == '"' {
				if err := buf.Seek(pos); err != nil {
					return nil, common.NewParseError(pos, err.Error())
				}

				return out, nil
			}

			if r == '\n' {
				return nil, common.NewParseError(pos, "newline in string literal")
			}

			if r == '\\' {
				e, err := parsecstrings.Any()(buf)
				if err != nil {
					return nil, common.NewParseError(pos, "unterminated escape")
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
func atEOL(buf common.Buffer[rune, parsecstrings.Position]) bool {
	return expr.AtEOL(buf)
}

// consumeEOL consumes '\n' if present.
func consumeEOL(buf common.Buffer[rune, parsecstrings.Position]) {
	if _, err := cNL(buf); err != nil {
		return // no '\n' - nothing to consume
	}
}

// consumeComment consumes a trailing comment if present (the combinator is
// passed by the backend); the absence of a comment is not an error.
func consumeComment(
	buf common.Buffer[rune, parsecstrings.Position],
	c common.Combinator[rune, parsecstrings.Position, string],
) {
	if _, err := c(buf); err != nil {
		return // no comment - nothing to consume
	}
}

// skipToEOL consumes everything up to end of line inclusive (recovery after
// a line parse error).
func skipToEOL(buf common.Buffer[rune, parsecstrings.Position]) {
	for {
		if _, err := cNotNL(buf); err != nil {
			break
		}
	}

	consumeEOL(buf)
}
