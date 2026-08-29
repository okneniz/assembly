// Package arb provides generators of arbitrary sets of mask rules and words
// (oh-snap Arbitrary) for dtree property tests. Masks of three kinds (sparse
// 1-4 bits, catch-all without bits, random broad), the match is either live
// (inside the mask) or dead (bits outside the mask - the rule is compatible
// with no word), up to 40 rules: decomposition into intersections, barriers,
// and bit duplicates is inevitable.
package arb

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"strings"

	ohsnap "github.com/okneniz/oh-snap"

	basearb "github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/dtree"
)

// Case is a Lookup property case: ordered rules and a word. Payload is the
// index of the rule, so that a mismatch shows WHOSE rule was chosen.
type Case struct {
	Rules []dtree.Rule[int]
	Word  uint32
}

func (c Case) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "word=%#08x rules=[", c.Word)
	for _, r := range c.Rules {
		fmt.Fprintf(&sb, "{mask:%#08x match:%#08x},", r.Mask, r.Match)
	}

	sb.WriteByte(']')
	return sb.String()
}

// LookupCase is an arbitrary case for the property "Lookup is equivalent to
// the linear list": Generate builds a rule set + a word, Shrink simplifies a
// failure (rules one at a time, the word toward zero, the masks toward zero).
func LookupCase(rnd *rand.Rand) ohsnap.Arbitrary[Case] {
	return lookupCase{rnd: rnd}
}

type lookupCase struct {
	rnd *rand.Rand
}

func (a lookupCase) Generate() iter.Seq[Case] {
	return basearb.Stream(func() Case {
		n := a.rnd.IntN(41)
		rules := make([]dtree.Rule[int], 0, n)

		for range n {
			var mask uint32
			switch a.rnd.IntN(4) {
			case 0:
				for range 1 + a.rnd.IntN(4) {
					mask |= 1 << a.rnd.IntN(32)
				}

			case 1: // catch-all
				mask = 0
			default:
				mask = a.rnd.Uint32()
			}

			match := a.rnd.Uint32() & mask
			if a.rnd.IntN(5) == 0 {
				match = a.rnd.Uint32() // dead rule
			}

			rules = append(rules, dtree.Rule[int]{
				Mask:    mask,
				Match:   match,
				Payload: len(rules),
			})
		}

		return Case{Rules: rules, Word: a.rnd.Uint32()}
	})
}

// Shrink produces simpler candidates: without each rule one at a time; the
// word toward zero (by clearing the lowest set bit); the mask of each rule
// toward zero.
func (lookupCase) Shrink(c Case) iter.Seq[Case] {
	var out []Case

	for i := range c.Rules {
		shorter := make([]dtree.Rule[int], 0, len(c.Rules)-1)
		shorter = append(shorter, c.Rules[:i]...)
		shorter = append(shorter, c.Rules[i+1:]...)
		out = append(out, Case{Rules: shorter, Word: c.Word})
	}

	if c.Word != 0 {
		out = append(out, Case{Rules: c.Rules, Word: c.Word & (c.Word - 1)})
	}

	for i, r := range c.Rules {
		if r.Mask == 0 {
			continue
		}

		zeroed := append([]dtree.Rule[int](nil), c.Rules...)
		zeroed[i] = dtree.Rule[int]{Payload: r.Payload}
		out = append(out, Case{Rules: zeroed, Word: c.Word})
	}

	return slices.Values(out)
}
