package asm

// The GAS macros: ".macro name[ param[=default]] [, param...]" captured
// raw until .endm, invoked by the bare name at the start of a line, the
// body lines expanded with \param substitution. An expansion re-enters
// the full source walk, so a macro may invoke a macro; the expanded
// lines keep their DEFINITION line numbers, so errors point into the
// macro body. v1: no label prefix before an invocation, a nested .macro
// inside a body refuses, the argument list is a plain comma split (a
// comment or the statement separator ends it first), and a missing
// argument without a default substitutes the empty string, as in GAS.

import (
	"fmt"

	parsec "github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

// macroParam is one formal: its name and the optional default.
type macroParam struct {
	name   string
	def    string
	hasDef bool
}

// macroLine is one raw body line with its source line number.
type macroLine struct {
	text []rune
	line int
}

// macroDef is one macro: the formals and the raw body.
type macroDef struct {
	name   string
	params []macroParam
	body   []macroLine
}

// parseMacroHeader is the .macro argument line: the name, then the
// parameters - comma-separated, the first one may follow the name
// without a comma (the seL4 spelling ".macro ventry label").
func parseMacroHeader(head []rune) (*macroDef, error) {
	parts := splitMacroComma(head)
	if len(parts) == 0 {
		return nil, fmt.Errorf(".macro: name expected")
	}

	words := splitMacroSpaces(parts[0])
	if len(words) == 0 || words[0] == "" {
		return nil, fmt.Errorf(".macro: name expected")
	}

	m := &macroDef{name: words[0]}
	paramWords := words[1:]
	paramWords = append(paramWords, splitFlat(parts[1:])...)
	for _, w := range paramWords {
		p, err := newMacroParam(w)
		if err != nil {
			return nil, err
		}

		m.params = append(m.params, p)
	}

	return m, nil
}

// newMacroParam is one parameter word: "name" or "name=default".
func newMacroParam(word string) (macroParam, error) {
	if word == "" {
		return macroParam{}, fmt.Errorf(".macro: parameter name expected")
	}

	for i := 0; i < len(word); i++ {
		if word[i] == '=' {
			return macroParam{
				name:   word[:i],
				def:    word[i+1:],
				hasDef: true,
			}, nil
		}
	}

	return macroParam{name: word}, nil
}

// splitMacroComma splits the raw text on commas (no nesting at this
// level); an empty text yields no parts.
func splitMacroComma(text []rune) [][]rune {
	var out [][]rune
	start := 0
	for i, r := range text {
		if r == ',' {
			out = append(out, text[start:i])
			start = i + 1
		}
	}

	return append(out, text[start:])
}

// splitMacroSpaces splits on runs of spaces and tabs.
func splitMacroSpaces(text []rune) []string {
	var out []string
	start := -1
	for i, r := range text {
		if r == ' ' || r == '\t' {
			if start >= 0 {
				out = append(out, string(text[start:i]))
				start = -1
			}

			continue
		}

		if start < 0 {
			start = i
		}
	}

	if start >= 0 {
		out = append(out, string(text[start:]))
	}

	return out
}

// splitFlat flattens the comma parts into parameter words.
func splitFlat(parts [][]rune) []string {
	var out []string
	for _, p := range parts {
		out = append(out, splitMacroSpaces(p)...)
	}

	return out
}

// macroArgs is the invocation argument list of the line after the name
// (the comment and the statement separator already cut); a blank line
// is no arguments.
func macroArgs(rest []rune) []string {
	if trimMacroSpaces(rest) == "" {
		return nil
	}

	parts := splitMacroComma(rest)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, trimMacroSpaces(p))
	}

	return out
}

// trimMacroSpaces cuts the leading and trailing spaces/tabs.
func trimMacroSpaces(text []rune) string {
	start, end := 0, len(text)
	for start < end && (text[start] == ' ' || text[start] == '\t') {
		start++
	}

	for end > start && (text[end-1] == ' ' || text[end-1] == '\t') {
		end--
	}

	return string(text[start:end])
}

// cutMacroLine cuts the raw .macro header line at the statement
// separator and the comment (quote-aware; the separator inside a string
// is data).
func cutMacroLine(text []rune, sep rune) []rune {
	inStr := false
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '"':
			inStr = !inStr
		case !inStr && sep != 0 && text[i] == sep:
			return text[:i]
		case !inStr && text[i] == '/' && i+1 < len(text) && text[i+1] == '/':
			return text[:i]
		}
	}

	return text
}

// readMacroArgs consumes the invocation argument text of the buffer: up
// to the statement separator (left in place), a '//' comment (consumed
// with the rest of the line), or the end of the line (its newline
// consumed) - quote-aware, a separator inside a string is data. atSep
// reports the separator stop.
func readMacroArgs(buf parsec.Buffer[rune, parsecstrings.Position], sep rune) (rest []rune, atSep bool) {
	var out []rune
	inStr := false
	for {
		r, ok := expr.PeekRune(buf)
		if !ok {
			return out, false
		}

		switch {
		case r == '\n':
			_ = expr.ConsumeRune(buf)
			return out, false
		case !inStr && sep != 0 && r == sep:
			return out, true
		case !inStr && r == '/':
			_ = expr.ConsumeRune(buf)
			next, nok := expr.PeekRune(buf)
			if nok && next == '/' {
				_ = readRawLine(buf) // the comment runs to the end of the line
				return out, false
			}

			out = append(out, r)
			continue
		case r == '"':
			inStr = !inStr
		}

		_ = expr.ConsumeRune(buf)
		out = append(out, r)
	}
}

// macroExpansion substitutes the arguments into the body and pads the
// lines with newlines so they keep their definition numbers (the walk
// of the expansion reports the body's source lines).
func macroExpansion(m *macroDef, args []string) ([]rune, error) {
	if len(args) > len(m.params) {
		return nil, fmt.Errorf(
			"macro %s: %d arguments for %d parameters",
			m.name,
			len(args),
			len(m.params),
		)
	}

	var sb []rune
	prev := 0
	for _, bl := range m.body {
		text, err := substituteMacro(bl.text, m, args)
		if err != nil {
			return nil, err
		}

		for range bl.line - prev - 1 {
			sb = append(sb, '\n')
		}

		sb = append(sb, text...)
		sb = append(sb, '\n')
		prev = bl.line
	}

	return sb, nil
}

// substituteMacro replaces \param with the actual value (or the
// default, or the empty string, as in GAS); at each backslash the
// longest parameter name wins (label before label2).
func substituteMacro(line []rune, m *macroDef, args []string) ([]rune, error) {
	var out []rune
	for i := 0; i < len(line); i++ {
		if line[i] != '\\' {
			out = append(out, line[i])
			continue
		}

		best := -1
		for pi, p := range m.params {
			if p.name != "" && runeHasPrefix(line[i+1:], p.name) &&
				(best < 0 || len(p.name) > len(m.params[best].name)) {
				best = pi
			}
		}

		if best < 0 {
			return nil, fmt.Errorf("macro %s: unknown parameter \\%s",
				m.name, macroWordAt(line[i+1:]))
		}

		if best < len(args) {
			out = append(out, []rune(args[best])...)
		} else {
			out = append(out, []rune(m.params[best].def)...)
		}

		i += len(m.params[best].name)
	}

	return out, nil
}

// runeHasPrefix is strings.HasPrefix over runes.
func runeHasPrefix(text []rune, prefix string) bool {
	pr := []rune(prefix)
	if len(pr) > len(text) {
		return false
	}

	for i, r := range pr {
		if text[i] != r {
			return false
		}
	}

	return true
}

// macroWordAt is the identifier at the start of the text (the unknown
// parameter name of the error).
func macroWordAt(text []rune) string {
	end := 0
	for end < len(text) && (expr.IsIdentStart(text[end]) || expr.IsIdentCont(text[end]) ||
		text[end] >= '0' && text[end] <= '9') {
		end++
	}

	return string(text[:end])
}

// readRawLine consumes the runes up to and including the newline,
// returning the text without it; at EOF what is left.
func readRawLine(buf parsec.Buffer[rune, parsecstrings.Position]) []rune {
	var out []rune
	for {
		r, ok := expr.PeekRune(buf)
		if !ok {
			return out
		}

		_ = expr.ConsumeRune(buf)
		if r == '\n' {
			return out
		}

		out = append(out, r)
	}
}

// firstMacroWord is the first space-separated word of a raw line (the
// .endm/.macro detection of the body capture).
func firstMacroWord(line []rune) string {
	for _, w := range splitMacroSpaces(line) {
		return w
	}

	return ""
}

// tryMacroWord consumes the word if it starts the buffer at a statement
// boundary (spaces, the separator, a comma, or the end of the line).
func tryMacroWord(buf parsec.Buffer[rune, parsecstrings.Position], word string) bool {
	save := buf.Position()
	for _, r := range word {
		p, ok := expr.PeekRune(buf)
		if !ok || p != r {
			_ = expr.Rewind(buf, save)
			return false
		}

		_ = expr.ConsumeRune(buf)
	}

	p, ok := expr.PeekRune(buf)
	if ok && p != ' ' && p != '\t' && p != '\n' && p != ',' {
		_ = expr.Rewind(buf, save)
		return false
	}

	return true
}

// scanMacroName consumes a leading identifier of the buffer.
func scanMacroName(buf parsec.Buffer[rune, parsecstrings.Position]) (string, bool) {
	var out []rune
	for {
		r, ok := expr.PeekRune(buf)
		if !ok || !(expr.IsIdentStart(r) || expr.IsIdentCont(r)) {
			return string(out), len(out) > 0
		}

		out = append(out, r)
		_ = expr.ConsumeRune(buf)
	}
}

// macroParseError is a parse error at a source position.
func macroParseError(pos parsecstrings.Position, format string, args ...any) parsec.Error[parsecstrings.Position] {
	return parsec.NewParseError(pos, fmt.Sprintf(format, args...))
}
