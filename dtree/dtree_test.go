package dtree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Cases: rules are {mask, match, payload}; a word; the expectation is the
// payload of the first matching rule of the list (or no match). The order of
// rules in the table = priority.

func TestLookupFirstMatch(t *testing.T) {
	cases := []struct {
		name  string
		rules []Rule[string]
		word  uint32
		want  string
		ok    bool
	}{
		{
			name: "exact match of a single rule",
			rules: []Rule[string]{
				{Mask: 0b1111_0000, Match: 0b1010_0000, Payload: "A"},
				{Mask: 0b1111_0000, Match: 0b0110_0000, Payload: "B"},
			},
			word: 0b1010_0101,
			want: "A",
			ok:   true,
		},
		{
			name: "intersection - the first in order wins",
			rules: []Rule[string]{
				{Mask: 0b1000, Match: 0b1000, Payload: "A"}, // catches 1xxx
				{Mask: 0b1100, Match: 0b1100, Payload: "B"}, // catches 11xx
			},
			word: 0b1100, // belongs to both
			want: "A",
			ok:   true,
		},
		{
			name: "masking entry - zero payload allowed",
			rules: []Rule[string]{
				{Mask: 0b1111_0000, Match: 0b0001_0000, Payload: ""},
				{Mask: 0b0000_0000, Match: 0b0000_0000, Payload: "catch-all"},
			},
			word: 0b0001_1010,
			want: "",
			ok:   true,
		},
		{
			name: "no match",
			rules: []Rule[string]{
				{Mask: 0b1111_0000, Match: 0b1010_0000, Payload: "A"},
			},
			word: 0b0101_0000,
			ok:   false,
		},
		{
			name: "duplicated bit - a rule without the bit in its mask is reachable in both branches",
			rules: []Rule[string]{
				{Mask: 0b0000, Match: 0b0000, Payload: "free"}, // indifferent to the high bits
				{Mask: 0b1000, Match: 0b1000, Payload: "hi"},
				{Mask: 0b0100, Match: 0b0100, Payload: "lo"},
			},
			word: 0b1000, // both free (first) and hi match
			want: "free",
			ok:   true,
		},
		{
			name: "masks in the high 32-bit bits",
			rules: []Rule[string]{
				{Mask: 0xF000_0000, Match: 0x8000_0000, Payload: "hi"},
				{Mask: 0xF000_0000, Match: 0x6000_0000, Payload: "lo"},
			},
			word: 0x8000_0042,
			want: "hi",
			ok:   true,
		},
		{
			name: "leaf list: word misses all rules of the leaf",
			rules: []Rule[string]{
				{Mask: 0b1000, Match: 0b1000, Payload: "A"}, // 1xxx
				{Mask: 0b1100, Match: 0b1100, Payload: "B"}, // 11xx
			},
			word: 0b0100, // neither 1xxx nor 11xx
			ok:   false,
		},
		{
			name: "node: field value without a branch",
			rules: []Rule[string]{
				{Mask: 0b1111_0000, Match: 0b1010_0000, Payload: "A"},
				{Mask: 0b1111_0000, Match: 0b0110_0000, Payload: "B"},
			},
			word: 0b0001_0000, // field value 00 - no branch
			ok:   false,
		},
		{
			name: "empty tree",
			word: 0xdeadbeef,
			ok:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := New(c.rules).Lookup(c.word)
			require.Equal(t, c.ok, ok)
			require.Equal(t, c.want, got)
		})
	}
}

// TestMetrics is the tree sizes: nodes/leaves/depth.
func TestMetrics(t *testing.T) {
	cases := []struct {
		name          string
		tree          *Tree[int]
		nodes, leaves int
		maxDepth      int
	}{
		{
			name:     "zero tree (without New)",
			tree:     &Tree[int]{},
			nodes:    0,
			leaves:   0,
			maxDepth: 0,
		},
		{
			name:     "single rule - root leaf",
			tree:     New([]Rule[int]{{Mask: 0xF0, Match: 0xA0, Payload: 1}}),
			nodes:    1,
			leaves:   1,
			maxDepth: 1,
		},
		{
			name: "vary-fork: root + two leaves",
			tree: New([]Rule[int]{
				{Mask: 0b1111_0000, Match: 0b1010_0000, Payload: 1},
				{Mask: 0b1111_0000, Match: 0b0110_0000, Payload: 2},
			}),
			nodes:    3,
			leaves:   2,
			maxDepth: 2,
		},
		{
			name: "intersection leaf list: a single root leaf",
			tree: New([]Rule[int]{
				{Mask: 0b1000, Match: 0b1000, Payload: 1},
				{Mask: 0b1100, Match: 0b1100, Payload: 2},
			}),
			nodes:    1,
			leaves:   1,
			maxDepth: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.nodes, c.tree.Nodes())
			require.Equal(t, c.leaves, c.tree.Leaves())
			require.Equal(t, c.maxDepth, c.tree.MaxDepth())
		})
	}
}

// TestLookupAgainstList is an oracle: over the enumeration of all 16-bit
// words, the tree must agree rule-for-rule with the linear list (first match).
func TestLookupAgainstList(t *testing.T) {
	rules := []Rule[int]{
		{Mask: 0b1100_0000_0000_0000, Match: 0b1000_0000_0000_0000, Payload: 1},
		{Mask: 0b1111_0000_0000_0000, Match: 0b0110_0000_0000_0000, Payload: 2},
		{Mask: 0b0000_1111_0000_0000, Match: 0b0000_0101_0000_0000, Payload: 3},
		{Mask: 0b0000_0000_1100_0000, Match: 0b0000_0000_1000_0000, Payload: 4},
		{Mask: 0b0000_0000_0000_0000, Match: 0b0000_0000_0000_0000, Payload: 5}, // catch-all
	}

	tree := New(rules)

	for w := range uint32(0x1_0000) {
		want, wantOK := 0, false
		for _, r := range rules {
			if w&r.Mask == r.Match {
				want, wantOK = r.Payload, true
				break
			}
		}

		got, gotOK := tree.Lookup(w)
		require.Equal(t, wantOK, gotOK, "word %#06x", w)
		require.Equal(t, want, got, "word %#06x", w)
	}
}
