package expr

// Expression lexing (parsec/strings): GAS integer literals (0x/0b/
// octal with a leading zero), character literals, symbol names, numeric
// local references. The base parsing rune atoms also live here -
// spaces, digits, name predicates, buffer peek/rewind: the expression is
// the lowest text layer of the assembler, the shared layer lives here, and
// the asm package builds the line grammar on it (asm imports expr, not the
// other way around).
//
// All atoms are Try-wrapped: a failed greedy atom leaves the position
// advanced, restoring it is the caller alternative's job (the parsec
// library contract).

import (
	"strconv"
	"unicode"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"
)

var (
	// cSpace is spaces EXCEPT newline: the assembler grammar is
	// line-based, '\n' terminates a statement (unicode.IsSpace('\n') ==
	// true, hence the explicit exclusion).
	cSpace = parsecstrings.Try(parsecstrings.Satisfy("whitespace", true, func(r rune) bool {
		return r != '\n' && unicode.IsSpace(r)
	}))
)

var (
	// CNL is the newline (line-based grammar).
	CNL = parsecstrings.Try(parsecstrings.Eq("newline", '\n'))
	// CNotNL is any character except newline.
	CNotNL = parsecstrings.Try(parsecstrings.NotEq("not a newline", '\n'))
	// CComma is the comma (operand separator).
	CComma = parsecstrings.Try(parsecstrings.Eq("comma", ','))

	cDecDigit = parsecstrings.Try(parsecstrings.Digit("decimal digit"))
	cHexDigit = parsecstrings.Try(parsecstrings.OneOf("hex digit",
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F'))
	cBinDigit = parsecstrings.Try(parsecstrings.OneOf("binary digit", '0', '1'))
	cOctDigit = parsecstrings.Try(
		parsecstrings.OneOf("octal digit", '0', '1', '2', '3', '4', '5', '6', '7'),
	)

	cPlus   = parsecstrings.Try(parsecstrings.Eq("'+'", '+'))
	cMinus  = parsecstrings.Try(parsecstrings.Eq("'-'", '-'))
	cTilde  = parsecstrings.Try(parsecstrings.Eq("'~'", '~'))
	cLParen = parsecstrings.Try(parsecstrings.Eq("'('", '('))
	cRParen = parsecstrings.Try(parsecstrings.Eq("')'", ')'))
	cQuote  = parsecstrings.Try(parsecstrings.Eq("quote", '\''))
	cEsc    = parsecstrings.Try(parsecstrings.Eq("escape", '\\'))
)

// IsIdentStart is the first rune of a symbol name: [._$a-zA-Z].
func IsIdentStart(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '.' || r == '$'
}

// IsIdentCont is a continuation rune of a symbol name: the start alphabet
// + digits.
func IsIdentCont(r rune) bool {
	return IsIdentStart(r) || r >= '0' && r <= '9'
}

// cSymbolExpr is a symbol name in an expression:
// [._$a-zA-Z][._$a-zA-Z0-9]* - captured maximally greedily (start +
// continuation, digits after the first letter are legal: "foo2").
var cSymbolExpr = parsecstrings.Cast(
	parsecstrings.Concat(8,
		parsecstrings.Some(4, "symbol name",
			parsecstrings.Try(parsecstrings.Satisfy("symbol start", true, IsIdentStart))),
		parsecstrings.Many(8,
			parsecstrings.Try(parsecstrings.Satisfy("symbol char", true, IsIdentCont))),
	),
	func(rs []rune) (*Expr, error) {
		return Sym(string(rs)), nil
	},
)

// cNumber is a GAS integer literal: 0x/0X hex, 0b/0B bin, a leading 0 is
// octal, "0" is zero, decimal without a leading zero. All alternatives are
// Try-wrapped: a failed branch must restore the position.
var cNumber = parsecstrings.Choice("number",
	parsecstrings.Try(numPrefixed("0x", 16)),
	parsecstrings.Try(numPrefixed("0X", 16)),
	parsecstrings.Try(numPrefixed("0b", 2)),
	parsecstrings.Try(numPrefixed("0B", 2)),
	parsecstrings.Try(octalC()),
	parsecstrings.Try(zeroC),
	parsecstrings.Try(decimalC),
)

// numPrefixed is a literal with a prefix ("0x", "0b") and >=1 digit of the
// corresponding base → int64.
func numPrefixed(prefix string, base int) parsec.Combinator[rune, parsecstrings.Position, *Expr] {
	digits := cHexDigit
	if base == 2 {
		digits = cBinDigit
	}

	return parsecstrings.Cast(
		parsecstrings.Skip(
			parsecstrings.String("prefix", prefix),
			parsecstrings.Some(8, "number digits", digits),
		),
		func(rs []rune) (*Expr, error) {
			v, err := strconv.ParseUint(string(rs), base, 64)
			if err != nil {
				return nil, err
			}

			return Num(int64(v)), nil
		},
	)
}

// octalC is '0' and >=1 octal digit ("010" = 8).
func octalC() parsec.Combinator[rune, parsecstrings.Position, *Expr] {
	return parsecstrings.Cast(
		parsecstrings.Concat(
			2,
			parsecstrings.Count(
				1,
				"leading zero",
				parsecstrings.Try(parsecstrings.Eq("zero", '0')),
			),
			parsecstrings.Some(4, "octal digits", cOctDigit),
		),
		func(rs []rune) (*Expr, error) {
			v, err := strconv.ParseInt(string(rs), 8, 64)
			if err != nil {
				return nil, err
			}

			return Num(v), nil
		},
	)
}

// zeroC is a single "0".
var zeroC = parsecstrings.Cast(
	parsecstrings.Count(1, "zero", parsecstrings.Try(parsecstrings.Eq("zero", '0'))),
	func(rs []rune) (*Expr, error) {
		return Num(0), nil
	},
)

// decimalC is decimal without a leading zero (first digit 1-9).
var decimalC = parsecstrings.Cast(
	parsecstrings.Concat(8,
		parsecstrings.Count(1, "non-zero digit",
			parsecstrings.Try(parsecstrings.Satisfy("digit 1-9", true, func(r rune) bool {
				return r >= '1' && r <= '9'
			}))),
		parsecstrings.Many(4, cDecDigit),
	),
	func(rs []rune) (*Expr, error) {
		v, err := strconv.ParseInt(string(rs), 10, 64)
		if err != nil {
			return nil, err
		}

		return Num(v), nil
	},
)

// cLocalRef is a reference to a numeric local label: [0-9]+ 'b'|'f' ("0b",
// "1f"). After b/f there must NOT follow an identifier-continuation
// character: "0b1010" is a binary literal (cNumber), not the reference "0b"
// with a tail. Builds a Sym with the full name; the assembler resolver
// distinguishes such names by the suffix.
var cLocalRef = func(buf parsec.Buffer[rune, parsecstrings.Position]) (*Expr, parsec.Error[parsecstrings.Position]) {
	save := buf.Position()

	var digits []rune
	for {
		r, err := cDecDigit(buf)
		if err != nil {
			break
		}

		digits = append(digits, r)
	}

	d, derr := parsecstrings.Try(parsecstrings.OneOf("'b' or 'f'", 'b', 'f'))(buf)
	if derr != nil || len(digits) == 0 {
		return nil, rewindOrError(buf, save)
	}

	if r, ok := peekRune(buf); ok && IsIdentCont(r) {
		return nil, rewindOrError(buf, save)
	}

	return Sym(string(digits) + string(d)), nil
}

// rewindOrError rewinds the buffer to the saved position; a Seek error
// becomes a parse error (as in Rewind).
func rewindOrError(
	buf parsec.Buffer[rune, parsecstrings.Position],
	save parsecstrings.Position,
) parsec.Error[parsecstrings.Position] {
	if err := buf.Seek(save); err != nil {
		return parsec.NewParseError(save, "rewind: "+err.Error())
	}

	return parsec.NewParseError(save, "local label reference expected")
}

// cCharLit is a character literal 'a' / '\n' / '\\' / '\0' → int64.
var cCharLit = parsecstrings.Cast(
	parsecstrings.Between(cQuote, parsecstrings.Try(escapeOrAny()), cQuote),
	func(rs []rune) (*Expr, error) {
		if len(rs) != 1 {
			return nil, strconv.ErrSyntax
		}

		return Num(int64(rs[0])), nil
	},
)

// escapeOrAny is one character: either an escape sequence (returns the
// character itself) or an arbitrary one.
func escapeOrAny() parsec.Combinator[rune, parsecstrings.Position, []rune] {
	return parsecstrings.Choice(
		"character",
		parsecstrings.Try(func() parsec.Combinator[rune, parsecstrings.Position, []rune] {
			return func(buf parsec.Buffer[rune, parsecstrings.Position]) ([]rune, parsec.Error[parsecstrings.Position]) {
				if _, err := cEsc(buf); err != nil {
					return nil, err
				}

				pos := buf.Position()
				r, err := parsecstrings.Any()(buf)
				if err != nil {
					return nil, err
				}

				switch r {
				case 'n':
					return []rune{'\n'}, nil
				case 't':
					return []rune{'\t'}, nil
				case 'r':
					return []rune{'\r'}, nil
				case '0':
					return []rune{0}, nil
				case '\\', '\'', '"':
					return []rune{r}, nil
				}

				return nil, parsec.NewParseError(pos, "unknown escape \\"+string(r))
			}
		}()),
		parsecstrings.Cast(
			parsecstrings.Any(),
			func(r rune) ([]rune, error) {
				return []rune{r}, nil
			},
		),
	)
}

// skipWS consumes spaces (including \r and \t).
func skipWS(buf parsec.Buffer[rune, parsecstrings.Position]) {
	for {
		if _, err := cSpace(buf); err != nil {
			return
		}
	}
}

// peekRune returns the next rune without consuming it; ok=false at EOF.
func peekRune(buf parsec.Buffer[rune, parsecstrings.Position]) (rune, bool) {
	if buf.IsEOF() {
		return 0, false
	}

	r, err := buf.Read(false)
	if err != nil {
		return 0, false
	}

	return r, true
}

// Rewind rewinds the buffer to the position saved before a trial branch of
// the grammar (from buf.Position()). A Seek error is also a parse error, as
// in parsec.Try.
func Rewind(
	buf parsec.Buffer[rune, parsecstrings.Position],
	save parsecstrings.Position,
) parsec.Error[parsecstrings.Position] {
	if err := buf.Seek(save); err != nil {
		return parsec.NewParseError(save, "rewind: "+err.Error())
	}

	return nil
}

// ConsumeRune consumes a rune just peeked by PeekRune. A read error after a
// successful peek is impossible; if it happens anyway - a parse error.
func ConsumeRune(
	buf parsec.Buffer[rune, parsecstrings.Position],
) parsec.Error[parsecstrings.Position] {
	if _, err := buf.Read(true); err != nil {
		return parsec.NewParseError(buf.Position(), "read: "+err.Error())
	}

	return nil
}

// --- exported for asm and arch backends (operand grammars) ---

var (
	// CSpace is a space, except newline.
	CSpace = cSpace
)

// CDecDigit is a decimal digit (backend numeric literals).
func CDecDigit() parsec.Combinator[rune, parsecstrings.Position, rune] {
	return cDecDigit
}

// SkipSpaces consumes spaces (except newline).
func SkipSpaces(buf parsec.Buffer[rune, parsecstrings.Position]) {
	skipWS(buf)
}

// PeekRune returns the next rune without consuming it; ok=false at EOF.
func PeekRune(buf parsec.Buffer[rune, parsecstrings.Position]) (rune, bool) {
	return peekRune(buf)
}

// SkipHash consumes a single '#' (the objdump/GAS immediate prefix) if
// present. A Try error only means "'#' is absent" - that is allowed.
func SkipHash(buf parsec.Buffer[rune, parsecstrings.Position]) {
	if _, err := parsecstrings.Try(parsecstrings.Eq("'#'", '#'))(buf); err != nil {
		return
	}
}

// AtEOL is true when the end of line follows (newline or EOF).
func AtEOL(buf parsec.Buffer[rune, parsecstrings.Position]) bool {
	r, ok := PeekRune(buf)
	return !ok || r == '\n'
}
