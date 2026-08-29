package arb

import (
	"math/rand/v2"
	"slices"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/dtree"
)

func TestLookupCaseGenerate(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 1))
	gen := LookupCase(rnd)

	empty, nonEmpty := 0, 0
	for range 1000 {
		c := ohsnap.First(gen.Generate())
		require.LessOrEqual(t, len(c.Rules), 40)

		if len(c.Rules) == 0 {
			empty++
		} else {
			nonEmpty++
		}
	}

	require.Positive(t, empty) // empty sets are part of the material
	require.Positive(t, nonEmpty)
}

func TestLookupCaseShrink(t *testing.T) {
	base := Case{
		Rules: []dtree.Rule[int]{
			{Mask: 0b1100, Match: 0b1000, Payload: 0},
			{Mask: 0b0110, Match: 0b0100, Payload: 1},
		},
		Word: 0b1010,
	}

	// Shrink order: without each rule one at a time, the word without the
	// lowest set bit, the mask of each rule (non-zero) toward zero.
	got := slices.Collect(LookupCase(nil).Shrink(base))
	require.Len(t, got, 5)
	require.Equal(t, Case{Rules: base.Rules[1:], Word: 0b1010}, got[0])
	require.Equal(t, Case{Rules: base.Rules[:1], Word: 0b1010}, got[1])
	require.Equal(t, Case{Rules: base.Rules, Word: 0b1000}, got[2])
	require.Equal(t, Case{
		Rules: []dtree.Rule[int]{{Payload: 0}, base.Rules[1]},
		Word:  0b1010,
	}, got[3])
	require.Equal(t, Case{
		Rules: []dtree.Rule[int]{base.Rules[0], {Payload: 1}},
		Word:  0b1010,
	}, got[4])
}

func TestCaseString(t *testing.T) {
	c := Case{
		Rules: []dtree.Rule[int]{
			{Mask: 0xC000, Match: 0x8000, Payload: 0},
		},
		Word: 0x0040,
	}

	require.Equal(t, "word=0x00000040 rules=[{mask:0x0000c000 match:0x00008000},]", c.String())
}

func TestLookupCaseShrinkEdges(t *testing.T) {
	// A zero word is not shrunk by bit clearing; zero masks are not shrunk
	// by zeroing - the only candidate: removing the single rule.
	base := Case{
		Rules: []dtree.Rule[int]{{Payload: 0}}, // mask=0, match=0
		Word:  0,
	}

	got := slices.Collect(LookupCase(nil).Shrink(base))
	require.Equal(t, []Case{{Rules: []dtree.Rule[int]{}, Word: 0}}, got)
}
