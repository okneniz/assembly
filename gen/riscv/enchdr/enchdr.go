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

	"github.com/okneniz/parsec/common"
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

// InsnDecl is a DECLARE_INSN line: a mnemonic linked to a pair of macros.
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

func newParsedLine(kind lineKind, macro Macro, decl InsnDecl) parsedLine {
	return parsedLine{
		kind:  kind,
		macro: macro,
		decl:  decl,
	}
}

// --- atomic combinators (Try wrappers are mandatory: a failed greedy atom
// leaves the position advanced; backtracking is the caller's responsibility). ---

var (
	cSpace    = strings.Try(strings.Space("whitespace"))
	cNewline  = strings.Try(strings.Eq("newline", '\n'))
	cNotNewl  = strings.Try(strings.NotEq("not a newline", '\n'))
	cComma    = strings.Try(strings.Eq("comma", ','))
	cLParen   = strings.Try(strings.Eq("'('", '('))
	cRParen   = strings.Try(strings.Eq("')'", ')'))
	cUndersc  = strings.Try(strings.Eq("'_'", '_'))
	cDefine   = strings.Try(strings.String("expected #define", "#define"))
	cDeclare  = strings.Try(strings.String("expected DECLARE_INSN", "DECLARE_INSN"))
	cSpaces1  = strings.Some(4, "expected whitespace", cSpace)
	cHexDigit = strings.Try(strings.OneOf("hex digit",
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	))
	cKind = strings.MapStrings("expected MATCH, MASK or CSR", map[string]lineKind{
		"MATCH": kindMatch,
		"MASK":  kindMask,
		"CSR":   kindCSR,
	})
	cHexValue = strings.Cast(
		strings.Skip(
			strings.String("expected 0x prefix", "0x"),
			strings.Some(8, "expected hex number", cHexDigit),
		),
		castUInt32,
	)
	// to end of line: everything except '\n', then an optional newline
	// (EOF without a trailing '\n' is not an error).
	cToEOL = strings.SkipMany(cNotNewl, strings.Optional(cNewline, rune(0)))
	cOther = strings.Cast(
		strings.SkipMany(cNotNewl, strings.Optional(cNewline, rune(0))),
		func(rune) (parsedLine, error) {
			return newParsedLine(kindOther, Macro{}, InsnDecl{}), nil
		},
	)
)

var lineKindPrefix = map[lineKind]string{
	kindMatch: "MATCH",
	kindMask:  "MASK",
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
func ident(what string, ok func(rune) bool) common.Combinator[rune, strings.Position, string] {
	return strings.Cast(
		strings.Some(16, what, strings.Try(strings.Satisfy(what, true, ok))),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
}

var cUpperIdent = ident("expected macro name", isUpperIdent)
var cLowerIdent = ident("expected instruction name", isLowerIdent)

// defineLineC is "#define" (MATCH|MASK|CSR)_NAME 0xVALUE [rest of line].
var defineLineC = func() common.Combinator[rune, strings.Position, parsedLine] {
	return func(buf common.Buffer[rune, strings.Position]) (parsedLine, common.Error[strings.Position]) {
		if _, err := cDefine(buf); err != nil {
			return parsedLine{}, err
		}

		if _, err := cSpaces1(buf); err != nil {
			return parsedLine{}, err
		}

		kind, err := cKind(buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := cUndersc(buf); err != nil {
			return parsedLine{}, err
		}

		name, err := cUpperIdent(buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := cSpaces1(buf); err != nil {
			return parsedLine{}, err
		}

		value, err := cHexValue(buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := cToEOL(buf); err != nil {
			return parsedLine{}, err
		}

		return newParsedLine(kind, NewMacro(name, value), InsnDecl{}), nil
	}
}()

// declLineC is DECLARE_INSN(name, MATCH_X, MASK_X) [rest of line].
var declLineC = func() common.Combinator[rune, strings.Position, parsedLine] {
	return func(buf common.Buffer[rune, strings.Position]) (parsedLine, common.Error[strings.Position]) {
		if _, err := cDeclare(buf); err != nil {
			return parsedLine{}, err
		}

		open := strings.SkipMany(cSpace, cLParen)
		if _, err := open(buf); err != nil {
			return parsedLine{}, err
		}

		insnName, err := strings.SkipMany(cSpace, cLowerIdent)(buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := strings.SkipMany(cSpace, cComma)(buf); err != nil {
			return parsedLine{}, err
		}

		matchName, err := strings.SkipMany(cSpace, cUpperIdent)(buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := strings.SkipMany(cSpace, cComma)(buf); err != nil {
			return parsedLine{}, err
		}

		maskName, err := strings.SkipMany(cSpace, cUpperIdent)(buf)
		if err != nil {
			return parsedLine{}, err
		}

		if _, err := strings.SkipMany(cSpace, cRParen)(buf); err != nil {
			return parsedLine{}, err
		}

		if _, err := cToEOL(buf); err != nil {
			return parsedLine{}, err
		}

		return newParsedLine(kindDecl, Macro{}, NewInsnDecl(insnName, matchName, maskName)), nil
	}
}()

// linesC is the file as a sequence of lines: recognized forms or skip.
var linesC = strings.Many(64, strings.Choice("unrecognized line",
	strings.Try(defineLineC),
	strings.Try(declLineC),
	cOther,
))

// Parse parses encoding.h text. Unrecognized lines are ignored; the result
// preserves appearance order for consumers' first-wins policies.
func Parse(data []rune) (Header, common.Error[strings.Position]) {
	lines, err := strings.Parse(data, linesC)
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
