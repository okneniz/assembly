// Package opcodes is a parsec grammar for the vendored loongarch-opcodes
// tables (arch/loong64/data/la-*.txt): machine-readable LoongArch
// instruction encoding lists compiled from the published ISA manuals by
// the loongson-community project.
//
// Each parseLine is the fixed encoding word (every operand field zeroed), the
// upstream mnemonic, the canonical operand parseFormat string, then optional
// @annotations - the official ISA mnemonic (@orig_name; the upstream
// tables rename some instructions), the manual operand order/shape
// (@orig_fmt), and informational tags:
//
//	00108000 add.d     DJK     @qemu
//	03000000 cu52i.d   DJSk12  @orig_name=lu52i.d @qemu
//	24000000 ldox4.w   DJSk14  @orig_name=ldptr.w @orig_fmt=DJSk14ps2
//	06482000 tlbclr    EMPTY
//
// The operand parseFormat notation (the upstream README ABNF): register slots
// D/J/K/A at fixed bit positions, then immediate slots - a signedness
// letter plus one or more index-width segments whose fields concatenate
// MSB-first (Sd5k16: [4:0] holds imm[20:16], [25:10] holds imm[15:0]).
// ParseFormat decodes the notation; Slots.Mask derives the encoding mask.
package opcodes

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// Entry is one instruction parseLine of a table.
type Entry struct {
	Word     uint32   // fixed encoding word, operand fields zeroed
	Name     string   // upstream mnemonic ("cu52i.d")
	Format   string   // canonical operand parseFormat ("DJSk12")
	OrigName string   // official ISA mnemonic, @orig_name ("" if none)
	OrigFmt  string   // manual operand order/shape, @orig_fmt ("" if none)
	Tags     []string // the remaining annotations, verbatim ("qemu", "rev=1p10")
}

// NewEntry - an instruction parseLine from the parsed fields.
func NewEntry(word uint32, parseName, parseFormat, origName, origFmt string, tags []string) Entry {
	return Entry{
		Word:     word,
		Name:     parseName,
		Format:   parseFormat,
		OrigName: origName,
		OrigFmt:  origFmt,
		Tags:     tags,
	}
}

// OfficialName - the mnemonic as the ISA manual and the vendor toolchains
// spell it (the upstream tables' renames restored).
func (e Entry) OfficialName() string {
	if e.OrigName != "" {
		return e.OrigName
	}

	return e.Name
}

// --- the table grammar. ---
//
// tableGrammar is the whole parseLine grammar of a loongarch-opcodes table:
// the atoms (Try wrappers are mandatory: a failed greedy atom leaves the
// position advanced; backtracking is the caller's responsibility) and
// the line combinators. Built once per Parse call.

type tableGrammar struct {
	// parseHSpace - a horizontal space: whitespace except the newline. Field
	// separators and trailing runs must never cross a line boundary
	// (strings.Space matches '\n', which would silently swallow the
	// following blank lines into the previous entry).
	parseHSpace    parsec.Combinator[rune, strings.Position, rune]
	parseNewline   parsec.Combinator[rune, strings.Position, rune]
	parseSpaces1   parsec.Combinator[rune, strings.Position, []rune]
	parseHexWord   parsec.Combinator[rune, strings.Position, uint32]
	parseName      parsec.Combinator[rune, strings.Position, string]
	parseFormat    parsec.Combinator[rune, strings.Position, string]
	parseAnnoField parsec.Combinator[rune, strings.Position, string]
	parseLine      parsec.Combinator[rune, strings.Position, Entry]
}

func makeTableGrammar() *tableGrammar {
	g := &tableGrammar{
		parseHSpace:  strings.Try(strings.Satisfy("horizontal space", true, isHSpaceRune)),
		parseNewline: strings.Try(strings.Eq("newline", '\n')),
	}

	g.parseSpaces1 = strings.Some(4, "expected whitespace", g.parseHSpace)

	at := strings.Try(strings.Eq("'@'", '@'))
	hexDigit := strings.Try(strings.OneOf("hex digit",
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	))
	nameRune := strings.Try(strings.Satisfy("mnemonic character", true, isNameRune))
	fmtRune := strings.Try(strings.Satisfy("format character", true, isFormatRune))
	annoRune := strings.Try(strings.Satisfy("annotation character", true, isAnnoRune))

	g.parseHexWord = strings.Cast(strings.Some(8, "expected hex word", hexDigit), castWord)
	g.parseName = strings.Cast(strings.Some(24, "expected mnemonic", nameRune), castRunes)
	g.parseFormat = strings.Cast(strings.Some(16, "expected operand format", fmtRune), castRunes)
	anno := strings.Cast(
		strings.Skip(at, strings.Some(40, "expected annotation", annoRune)),
		castRunes,
	)

	// parseAnnoField - one annotation, separated from the previous field by
	// whitespace (Try rolls the spaces back when no '@' follows).
	g.parseAnnoField = strings.Try(strings.Skip(g.parseSpaces1, anno))

	g.parseLine = g.makeLineParser()

	return g
}

func isNameRune(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '.' || r == '_'
}

func isHSpaceRune(r rune) bool {
	return r != '\n' && unicode.IsSpace(r)
}

func isFormatRune(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

func isAnnoRune(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' ||
		r >= '0' && r <= '9' || r == '.' || r == '_' || r == '=' || r == '+'
}

func castRunes(rs []rune) (string, error) {
	return string(rs), nil
}

// castWord - the encoding word: exactly 8 hex digits (zero-padded upstream).
func castWord(ds []rune) (uint32, error) {
	if len(ds) != 8 {
		return 0, fmt.Errorf("hex word: want 8 digits, got %d", len(ds))
	}

	v, err := strconv.ParseUint(string(ds), 16, 32)
	if err != nil {
		return 0, fmt.Errorf("hex word: %w", err)
	}

	return uint32(v), nil
}

// newEntryAssembled - an Entry from the raw parseLine fields, with the
// annotations classified (orig_name/orig_fmt extracted, the rest are tags).
func newEntryAssembled(word uint32, parseName, parseFormat string, annos []string) Entry {
	origName, origFmt := "", ""
	var tags []string

	for _, a := range annos {
		if v, ok := cutPrefix(a, "orig_name="); ok {
			origName = v

			continue
		}

		if v, ok := cutPrefix(a, "orig_fmt="); ok {
			origFmt = v

			continue
		}

		tags = append(tags, a)
	}

	return NewEntry(word, parseName, parseFormat, origName, origFmt, tags)
}

// cutPrefix - s without the prefix and ok (a local strings.CutPrefix: the
// parsec package is imported under the stdlib strings parseName).
func cutPrefix(s, prefix string) (string, bool) {
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):], true
	}

	return s, false
}

// skipSpaces - consume the run of horizontal spaces.
func (g *tableGrammar) skipSpaces(buf parsec.Buffer[rune, strings.Position]) {
	for {
		if _, err := g.parseHSpace(buf); err != nil {
			return
		}
	}
}

// consumeNewline - the parseLine's parseNewline, if present.
func (g *tableGrammar) consumeNewline(buf parsec.Buffer[rune, strings.Position]) {
	if _, err := g.parseNewline(buf); err != nil {
		return
	}
}

// makeLineParser - one table line: word, mnemonic, format, annotations,
// EOL - or a blank line (whitespace only, the zero Entry which Parse
// drops).
func (g *tableGrammar) makeLineParser() parsec.Combinator[rune, strings.Position, Entry] {
	blankLine := strings.Try(
		func(buf parsec.Buffer[rune, strings.Position]) (Entry, parsec.Error[strings.Position]) {
			if _, err := g.parseHSpace(buf); err != nil {
				if _, err := g.parseNewline(buf); err != nil {
					return Entry{}, parsec.NewParseError(
						buf.Position(),
						"expected a blank line",
					)
				}

				return Entry{}, nil
			}

			g.skipSpaces(buf)
			g.consumeNewline(buf)

			return Entry{}, nil
		},
	)

	tableLine := func(buf parsec.Buffer[rune, strings.Position]) (Entry, parsec.Error[strings.Position]) {
		word, err := g.parseHexWord(buf)
		if err != nil {
			return Entry{}, err
		}

		if _, err := g.parseSpaces1(buf); err != nil {
			return Entry{}, err
		}

		name, err := g.parseName(buf)
		if err != nil {
			return Entry{}, err
		}

		if _, err := g.parseSpaces1(buf); err != nil {
			return Entry{}, err
		}

		format, err := g.parseFormat(buf)
		if err != nil {
			return Entry{}, err
		}

		// The annotations: '@' tokens separated by whitespace. The
		// annoField is a Try - when no '@' follows the spaces it fails
		// with the position rolled back, and the loop simply ends.
		annos := []string{}
		for a, aerr := g.parseAnnoField(buf); aerr == nil; a, aerr = g.parseAnnoField(buf) {
			annos = append(annos, a)
		}

		g.skipSpaces(buf)
		g.consumeNewline(buf)

		return newEntryAssembled(word, name, format, annos), nil
	}

	return strings.Choice("unrecognized line",
		strings.Try(tableLine),
		blankLine,
	)
}

// Parse parses one loongarch-opcodes table in file order. Every non-blank
// parseLine must be a full entry: unlike encoding.h there are no comments, so
// a parseLine matching neither form is corrupt data and an error.
func Parse(data []rune) ([]Entry, parsec.Error[strings.Position]) {
	g := makeTableGrammar()
	buf := strings.Buffer(data)

	out := make([]Entry, 0, 64)
	for !buf.IsEOF() {
		e, err := g.parseLine(buf)
		if err != nil {
			return nil, parsec.NewParseError(buf.Position(), "corrupt line", err)
		}

		if e.Name == "" {
			continue // blank parseLine
		}

		out = append(out, e)
	}

	return out, nil
}
