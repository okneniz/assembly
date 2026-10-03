// Package enchdr is a parsec grammar for the canonical Spike encoding.h
// header (arch/riscv/data/encoding.h). Shared by gen-riscv-instr
// (MATCH_*/MASK_* macros + DECLARE_INSN lines) and gen-riscv-csr
// (CSR_* definitions).
//
// The grammar is scanning: unrecognized lines (comments, #if/#endif,
// #define with _RISCV_UL(...) wrappers) are skipped to end of line, as the
// former regex approach did — their content is not decode data.
// Recognized forms (in file order):
//
//	#define MATCH_ADDI 0x13          → Macro{"MATCH_ADDI", 0x13}
//	#define MASK_ADDI 0x707f         → Macro{"MASK_ADDI", 0x707f}
//	#define CSR_FFLAGS 0x1           → Macro{"FFLAGS", 0x1} (no prefix)
//	DECLARE_INSN(addi, MATCH_ADDI, MASK_ADDI) → InsnDecl{...}
package enchdr

import (
	"fmt"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// Macro is a named encoding.h constant: MATCH_*/MASK_* (full name with
// prefix) or CSR_* (name without prefix, value is a 12-bit address).
type Macro struct {
	Name  string
	Value uint32
}

func NewMacro(name string, value uint32) Macro {
	return Macro{
		Name:  name,
		Value: value,
	}
}

// InsnDecl is a DECLARE_INSN parseLine: a mnemonic linked to a pair of macros.
type InsnDecl struct {
	Name      string // "addi"
	MatchName string // "MATCH_ADDI"
	MaskName  string // "MASK_ADDI"
}

func NewInsnDecl(name string, matchName string, maskName string) InsnDecl {
	return InsnDecl{
		Name:      name,
		MatchName: matchName,
		MaskName:  maskName,
	}
}

// Header is everything recognized in the header, in file order.
// Deduplication policies (first-wins by name/address) are applied by consumers.
type Header struct {
	Macros []Macro
	CSRs   []Macro
	Insns  []InsnDecl
}

type lineKind int

const (
	kindOther lineKind = iota
	kindMatch
	kindMask
	kindCSR
	kindDecl
)

type parsedLine struct {
	kind  lineKind
	macro Macro
	decl  InsnDecl
}

var lineKindPrefix = map[lineKind]string{
	kindMatch: "MATCH",
	kindMask:  "MASK",
}

// Parse parses encoding.h text. Unrecognized lines are ignored; the result
// preserves appearance order for consumers' first-wins policies.
func Parse(data []rune) (Header, parsec.Error[strings.Position]) {
	lines, err := strings.Parse(parsec.Stateless{}, data, makeLinesParser())
	if err != nil {
		return Header{}, err
	}

	h := Header{}
	for _, l := range lines {
		switch l.kind {
		case kindMatch, kindMask:
			h.Macros = append(
				h.Macros,
				NewMacro(lineKindPrefix[l.kind]+"_"+l.macro.Name, l.macro.Value),
			)
		case kindCSR:
			h.CSRs = append(h.CSRs, l.macro)
		case kindDecl:
			h.Insns = append(h.Insns, l.decl)
		case kindOther:
			// not a macro/declaration — skipped
		}
	}

	return h, nil
}

func newParsedLine(kind lineKind, macro Macro, decl InsnDecl) parsedLine {
	return parsedLine{
		kind:  kind,
		macro: macro,
		decl:  decl,
	}
}

// newLinesC builds the file grammar: a sequence of lines, each a
// recognized form (define/declaration) or a skipped other line. The
// combinators are values captured by the closures (Try wrappers are
// mandatory: a failed greedy atom leaves the position advanced;
// backtracking is the caller's responsibility).
func makeLinesParser() parsec.Combinator[rune, strings.Position, []parsedLine, parsec.Stateless] {
	space := strings.Try(strings.Space[parsec.Stateless]("whitespace"))
	newline := strings.Try(strings.Eq[parsec.Stateless]("newline", '\n'))
	notNewl := strings.Try(strings.NotEq[parsec.Stateless]("not a newline", '\n'))
	comma := strings.Try(strings.Eq[parsec.Stateless]("comma", ','))
	lparen := strings.Try(strings.Eq[parsec.Stateless]("'('", '('))
	rparen := strings.Try(strings.Eq[parsec.Stateless]("')'", ')'))
	undersc := strings.Try(strings.Eq[parsec.Stateless]("'_'", '_'))
	define := strings.Try(strings.String[parsec.Stateless]("expected #define", "#define"))
	declare := strings.Try(
		strings.String[parsec.Stateless]("expected DECLARE_INSN", "DECLARE_INSN"),
	)
	spaces1 := strings.Some(4, "expected whitespace", space)
	hexDigit := strings.Try(strings.OneOf[parsec.Stateless]("hex digit",
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	))
	kind := strings.MapStrings[lineKind, parsec.Stateless](
		"expected MATCH, MASK or CSR",
		map[string]lineKind{
			"MATCH": kindMatch,
			"MASK":  kindMask,
			"CSR":   kindCSR,
		},
	)
	hexValue := strings.Cast(
		strings.Skip(
			strings.String[parsec.Stateless]("expected 0x prefix", "0x"),
			strings.Some(8, "expected hex number", hexDigit),
		),
		castUInt32,
	)
	// to end of parseLine: everything except '\n', then an optional newline
	// (EOF without a trailing '\n' is not an error).
	toEOL := strings.SkipMany(notNewl, strings.Optional(newline, rune(0)))
	other := strings.Cast(
		strings.SkipMany(notNewl, strings.Optional(newline, rune(0))),
		func(rune) (parsedLine, error) {
			return newParsedLine(kindOther, Macro{}, InsnDecl{}), nil
		},
	)

	upperIdent := ident("expected macro name", isUpperIdent)
	lowerIdent := ident("expected instruction name", isLowerIdent)

	// defineLine is "#define" (MATCH|MASK|CSR)_NAME 0xVALUE [rest of line].
	defineLine := func(state parsec.Stateless, buf parsec.Buffer[rune, strings.Position]) (parsedLine, parsec.Error[strings.Position]) {
		if _, err := define(state, buf); err != nil {
			return parsedLine{}, err
		}

		if _, err := spaces1(state, buf); err != nil {
			return parsedLine{}, err
		}

		k, err := kind(state, buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := undersc(state, buf); err != nil {
			return parsedLine{}, err
		}

		name, err := upperIdent(state, buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := spaces1(state, buf); err != nil {
			return parsedLine{}, err
		}

		value, err := hexValue(state, buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := toEOL(state, buf); err != nil {
			return parsedLine{}, err
		}

		return newParsedLine(k, NewMacro(name, value), InsnDecl{}), nil
	}

	// declLine is DECLARE_INSN(name, MATCH_X, MASK_X) [rest of line].
	declLine := func(state parsec.Stateless, buf parsec.Buffer[rune, strings.Position]) (parsedLine, parsec.Error[strings.Position]) {
		if _, err := declare(state, buf); err != nil {
			return parsedLine{}, err
		}

		open := strings.SkipMany(space, lparen)
		if _, err := open(state, buf); err != nil {
			return parsedLine{}, err
		}

		insnName, err := strings.SkipMany(space, lowerIdent)(state, buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := strings.SkipMany(space, comma)(state, buf); err != nil {
			return parsedLine{}, err
		}

		matchName, err := strings.SkipMany(space, upperIdent)(state, buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := strings.SkipMany(space, comma)(state, buf); err != nil {
			return parsedLine{}, err
		}

		maskName, err := strings.SkipMany(space, upperIdent)(state, buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := strings.SkipMany(space, rparen)(state, buf); err != nil {
			return parsedLine{}, err
		}

		if _, err := toEOL(state, buf); err != nil {
			return parsedLine{}, err
		}

		return newParsedLine(kindDecl, Macro{}, NewInsnDecl(insnName, matchName, maskName)), nil
	}

	return strings.Many(64, strings.Choice("unrecognized line",
		strings.Try(defineLine),
		strings.Try(declLine),
		other,
	))
}

func castUInt32(ds []rune) (uint32, error) {
	v, err := strconv.ParseUint(string(ds), 16, 32)
	if err != nil {
		return 0, fmt.Errorf("hex value: %w", err)
	}

	return uint32(v), nil
}

func isUpperIdent(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
}

func isLowerIdent(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_'
}

// ident is a non-empty sequence of runes satisfying ok.
func ident(
	what string,
	ok func(rune) bool,
) parsec.Combinator[rune, strings.Position, string, parsec.Stateless] {
	return strings.Cast(
		strings.Some(16, what, strings.Try(strings.Satisfy[parsec.Stateless](what, true, ok))),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
}
