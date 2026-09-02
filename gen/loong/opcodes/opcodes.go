// Package opcodes is a parsec grammar for the vendored loongarch-opcodes
// tables (arch/loong64/data/la-*.txt): machine-readable LoongArch
// instruction encoding lists compiled from the published ISA manuals by
// the loongson-community project.
//
// Each line is the fixed encoding word (every operand field zeroed), the
// upstream mnemonic, the canonical operand format string, then optional
// @annotations - the official ISA mnemonic (@orig_name; the upstream
// tables rename some instructions), the manual operand order/shape
// (@orig_fmt), and informational tags:
//
//	00108000 add.d     DJK     @qemu
//	03000000 cu52i.d   DJSk12  @orig_name=lu52i.d @qemu
//	24000000 ldox4.w   DJSk14  @orig_name=ldptr.w @orig_fmt=DJSk14ps2
//	06482000 tlbclr    EMPTY
//
// The operand format notation (the upstream README ABNF): register slots
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

// Entry is one instruction line of a table.
type Entry struct {
	Word     uint32   // fixed encoding word, operand fields zeroed
	Name     string   // upstream mnemonic ("cu52i.d")
	Format   string   // canonical operand format ("DJSk12")
	OrigName string   // official ISA mnemonic, @orig_name ("" if none)
	OrigFmt  string   // manual operand order/shape, @orig_fmt ("" if none)
	Tags     []string // the remaining annotations, verbatim ("qemu", "rev=1p10")
}

// NewEntry - an instruction line from the parsed fields.
func NewEntry(word uint32, name, format, origName, origFmt string, tags []string) Entry {
	return Entry{
		Word:     word,
		Name:     name,
		Format:   format,
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

// --- atomic combinators (Try wrappers are mandatory: a failed greedy atom
// leaves the position advanced; backtracking is the caller's responsibility). ---

var (
	// cHSpace - a horizontal space: whitespace except the newline. Field
	// separators and trailing runs must never cross a line boundary
	// (strings.Space matches '\n', which would silently swallow the
	// following blank lines into the previous entry).
	cHSpace   = strings.Try(strings.Satisfy("horizontal space", true, isHSpaceRune))
	cNewline  = strings.Try(strings.Eq("newline", '\n'))
	cAt       = strings.Try(strings.Eq("'@'", '@'))
	cSpaces1  = strings.Some(4, "expected whitespace", cHSpace)
	cHexDigit = strings.Try(strings.OneOf("hex digit",
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	))
	cNameRune = strings.Try(strings.Satisfy("mnemonic character", true, isNameRune))
	cFmtRune  = strings.Try(strings.Satisfy("format character", true, isFormatRune))
	cAnnoRune = strings.Try(strings.Satisfy("annotation character", true, isAnnoRune))

	cHexWord = strings.Cast(strings.Some(8, "expected hex word", cHexDigit), castWord)
	cName    = strings.Cast(strings.Some(24, "expected mnemonic", cNameRune), castRunes)
	cFormat  = strings.Cast(strings.Some(16, "expected operand format", cFmtRune), castRunes)
	cAnno    = strings.Cast(
		strings.Skip(cAt, strings.Some(40, "expected annotation", cAnnoRune)),
		castRunes,
	)

	// cAnnoField - one annotation, separated from the previous field by
	// whitespace (Try rolls the spaces back when no '@' follows).
	cAnnoField = strings.Try(strings.Skip(cSpaces1, cAnno))
)

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

// newEntryAssembled - an Entry from the raw line fields, with the
// annotations classified (orig_name/orig_fmt extracted, the rest are tags).
func newEntryAssembled(word uint32, name, format string, annos []string) Entry {
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

	return NewEntry(word, name, format, origName, origFmt, tags)
}

// cutPrefix - s without the prefix and ok (a local strings.CutPrefix: the
// parsec package is imported under the stdlib strings name).
func cutPrefix(s, prefix string) (string, bool) {
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):], true
	}

	return s, false
}

// skipSpaces - consume the run of horizontal spaces.
func skipSpaces(buf parsec.Buffer[rune, strings.Position]) {
	for {
		if _, err := cHSpace(buf); err != nil {
			return
		}
	}
}

// consumeNewline - the line's newline, if present.
func consumeNewline(buf parsec.Buffer[rune, strings.Position]) {
	if _, err := cNewline(buf); err != nil {
		return
	}
}

// lineC - one table line: word, mnemonic, format, annotations, EOL.
var lineC = func() parsec.Combinator[rune, strings.Position, Entry] {
	return func(buf parsec.Buffer[rune, strings.Position]) (Entry, parsec.Error[strings.Position]) {
		word, err := cHexWord(buf)
		if err != nil {
			return Entry{}, err
		}

		if _, err := cSpaces1(buf); err != nil {
			return Entry{}, err
		}

		name, err := cName(buf)
		if err != nil {
			return Entry{}, err
		}

		if _, err := cSpaces1(buf); err != nil {
			return Entry{}, err
		}

		format, err := cFormat(buf)
		if err != nil {
			return Entry{}, err
		}

		// The annotations: '@' tokens separated by whitespace. cAnnoField
		// is a Try - when no '@' follows the spaces it fails with the
		// position rolled back, and the loop simply ends.
		annos := []string{}
		for {
			a, err := cAnnoField(buf)
			if err != nil {
				break
			}

			annos = append(annos, a)
		}

		skipSpaces(buf)
		consumeNewline(buf)

		return newEntryAssembled(word, name, format, annos), nil
	}
}()

// cBlankLine - a whitespace-only line (horizontal spaces and/or one
// newline; never zero-width), as the zero Entry (Name == ""), which Parse
// drops.
var cBlankLine = strings.Try(func() parsec.Combinator[rune, strings.Position, Entry] {
	return func(buf parsec.Buffer[rune, strings.Position]) (Entry, parsec.Error[strings.Position]) {
		if _, err := cHSpace(buf); err != nil {
			if _, err := cNewline(buf); err != nil {
				return Entry{}, parsec.NewParseError(buf.Position(), "expected a blank line")
			}

			return Entry{}, nil
		}

		skipSpaces(buf)
		consumeNewline(buf)

		return Entry{}, nil
	}
}())

// fileLine - a table line or a blank line.
var fileLine = strings.Choice("unrecognized line",
	strings.Try(lineC),
	cBlankLine,
)

// Parse parses one loongarch-opcodes table in file order. Every non-blank
// line must be a full entry: unlike encoding.h there are no comments, so
// a line matching neither form is corrupt data and an error.
func Parse(data []rune) ([]Entry, parsec.Error[strings.Position]) {
	buf := strings.Buffer(data)

	out := make([]Entry, 0, 64)
	for !buf.IsEOF() {
		e, err := fileLine(buf)
		if err != nil {
			return nil, parsec.NewParseError(buf.Position(), "corrupt line", err)
		}

		if e.Name == "" {
			continue // blank line
		}

		out = append(out, e)
	}

	return out, nil
}
