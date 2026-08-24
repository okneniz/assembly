package main

// parsec grammar for extracting the mnemonic from the asmtemplate text of the
// A64 ISA XML. The previous approach (strings.Fields + a per-character loop)
// looked for the first alphabetic run of the first token; since tokens are
// separated by spaces, this is equivalent to the first maximal alphabetic
// run of the whole text ("ADD", "NOP", "BL"), lowercased.

import (
	"fmt"
	"strings"

	parsecstrings "github.com/okneniz/parsec/strings"
)

func isLatin(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

var (
	cNotAlpha = parsecstrings.Try(parsecstrings.Satisfy("non-letter", true, func(r rune) bool {
		return !isLatin(r)
	}))
	cAlphaRun = parsecstrings.Cast(
		parsecstrings.Some(8, "mnemonic letters",
			parsecstrings.Try(parsecstrings.Satisfy("letter", true, isLatin))),
		func(rs []rune) (string, error) {
			return strings.ToLower(string(rs)), nil
		},
	)
	// cMnemonic — skip all non-letters, then read the first alphabetic
	// run; if there are no letters at all — an empty string (as in the
	// previous mnemonic).
	cMnemonic = parsecstrings.Optional(
		parsecstrings.SkipMany(cNotAlpha, cAlphaRun),
		"",
	)
)

// firstAlphaRun returns the first maximal alphabetic run of the text
// (asmtemplate chunks are joined with a space), lowercased; "" if there are
// no letters.
func firstAlphaRun(pieces []string) (string, error) {
	var buf []rune
	for _, p := range pieces {
		buf = append(buf, []rune(p)...)
		buf = append(buf, ' ')
	}

	m, err := parsecstrings.Parse(buf, cMnemonic)
	if err != nil {
		// cMnemonic is Optional with a default "", so an error can only
		// come from the buffer; the generator's caller sees it upstream
		return "", fmt.Errorf("firstAlphaRun: %w", err)
	}

	return m, nil
}
