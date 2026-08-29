package expr

// Junk string generator: parser robustness — arbitrary junk from the runes
// of the expression lexicon either parses or is rejected with an error.
// Not ohsnap.ArbitraryString: that one picks runes with a global rand.IntN
// (not the passed rnd), the seed is not reproduced — hence a custom
// generator for deterministic failures.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
)

// junkAlphabet — runes of the expression lexicon: digits and base prefixes,
// operators, brackets, quotes and escapes of character literals, separators, name characters.
const junkAlphabet = "0123456789abcdefxXbB+-*/%&|^~<>() '\"\\._$\t\r"

// junkGen — junk string generator, length 0..32.
type junkGen struct {
	rnd *rand.Rand
}

// Junk — an arbitrary string from the runes of the expression lexicon.
func Junk(rnd *rand.Rand) ohsnap.Arbitrary[string] {
	return junkGen{rnd: rnd}
}

func (g junkGen) Generate() iter.Seq[string] {
	return arb.Stream(func() string {
		rs := make([]rune, g.rnd.IntN(33))
		for i := range rs {
			rs[i] = rune(junkAlphabet[g.rnd.IntN(len(junkAlphabet))])
		}

		return string(rs)
	})
}

// Shrink — prefix-half and "drop one" (like arb.Seq): a minimal
// failing substring localizes the bug.
func (junkGen) Shrink(s string) iter.Seq[string] {
	r := []rune(s)
	var out []string
	if h := len(r) / 2; h > 0 {
		out = append(out, string(r[:h]))
	}

	for i := range r {
		out = append(out, string(append(append([]rune{}, r[:i]...), r[i+1:]...)))
	}

	return slices.Values(out)
}
