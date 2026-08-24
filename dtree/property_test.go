package dtree_test

import (
	mrnd "math/rand/v2"
	"os"
	"strconv"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/dtree"
	"github.com/okneniz/assembly/dtree/arb"
)

// propFirstMatch is an oracle: the first match of the linear list.
func propFirstMatch(c arb.Case) (payload int, ok bool) {
	for _, r := range c.Rules {
		if c.Word&r.Mask == r.Match {
			return r.Payload, true
		}
	}

	return 0, false
}

// TestPropertyLookupMatchesList: Lookup is equivalent to the linear list on
// any rules and words: it returns the payload of the first (in New order)
// matching rule. The generator is dtree/arb (three kinds of masks, dead
// matches, catch-alls), the seed is ASSEMBLY_SEED (default 42).
func TestPropertyLookupMatchesList(t *testing.T) {
	seed := uint64(42)
	if s := os.Getenv("ASSEMBLY_SEED"); s != "" {
		if v, err := strconv.ParseUint(s, 0, 64); err == nil {
			seed = v
		}
	}

	t.Logf("seed: %d (ASSEMBLY_SEED)", seed)
	cases := arb.LookupCase(mrnd.New(mrnd.NewPCG(seed, seed)))

	ohsnap.Check(t, 100000, cases, func(c arb.Case) bool {
		want, wantOK := propFirstMatch(c)
		got, gotOK := dtree.New(c.Rules).Lookup(c.Word)

		if wantOK != gotOK {
			t.Errorf("ok: list=%v tree=%v (%s)", wantOK, gotOK, c)
			return false
		}

		if wantOK && want != got {
			t.Errorf("payload: list=%d tree=%d (%s)", want, got, c)
			return false
		}

		return true
	})
}

// propSeed is the seed from ASSEMBLY_SEED (default 42), logged.
func propSeed(tb testing.TB) *mrnd.Rand {
	tb.Helper()

	seed := uint64(42)
	if s := os.Getenv("ASSEMBLY_SEED"); s != "" {
		if v, err := strconv.ParseUint(s, 0, 64); err == nil {
			seed = v
		}
	}

	tb.Logf("seed: %d (ASSEMBLY_SEED)", seed)

	return mrnd.New(mrnd.NewPCG(seed, seed))
}

// TestPropertyDepthBoundedByWordBits: the depth of the tree is bounded by the
// number of word bits (32), not by the number of rules: the cost of Lookup
// does not grow with the size of the registry.
func TestPropertyDepthBoundedByWordBits(t *testing.T) {
	cases := arb.LookupCase(propSeed(t))

	ohsnap.Check(t, 100000, cases, func(c arb.Case) bool {
		if d := dtree.New(c.Rules).MaxDepth(); d > 32 {
			t.Errorf("depth %d > 32 word bits (%s)", d, c)

			return false
		}

		return true
	})
}

// TestPropertyDuplicatesCollapse: K copies of a single rule land in a single
// leaf list: the size of the tree does not grow with the number of
// duplicates - reuse of a path that the linear list does not have. (The
// strong form "Stats does not change when adding a copy to an arbitrary set"
// is false: a copy changes the len of the group, and with it the duplication
// threshold and the choice of the cutting bit.)
func TestPropertyDuplicatesCollapse(t *testing.T) {
	protos := []dtree.Rule[int]{
		{Mask: 0b1111_0000, Match: 0b1010_0000, Payload: 0},
		{Mask: 0x8000_0000, Match: 0x8000_0000, Payload: 0},
		{Mask: 0xF000_0F00, Match: 0x6000_0500, Payload: 0},
		{Mask: 0, Match: 0, Payload: 0}, // catch-all
	}

	for _, p := range protos {
		for _, k := range []int{1, 2, 3, 8, 33, 100} {
			rules := make([]dtree.Rule[int], k)
			for i := range rules {
				rules[i] = p
			}

			tree := dtree.New(rules)
			require.Equal(t, 1, tree.Nodes(), "proto=%#v k=%d", p, k)
			require.Equal(t, 1, tree.Leaves(), "proto=%#v k=%d", p, k)
			require.Equal(t, 1, tree.MaxDepth(), "proto=%#v k=%d", p, k)
		}
	}
}

// TestPropertyRefinementFamilySingleLeaf: a hierarchy of refinement rules
// (the mask grows bit by bit, the match is shared) lands in a SINGLE leaf
// list: the size of the tree does not grow with the number of refinements -
// the sublinearity the tree exists for (unlike the linear list).
func TestPropertyRefinementFamilySingleLeaf(t *testing.T) {
	for _, k := range []int{1, 2, 3, 5, 8, 16, 33, 64} {
		t.Run(strconv.Itoa(k), func(t *testing.T) {
			rules := make([]dtree.Rule[int], 0, k)

			var mask uint32
			for b := uint(0); len(rules) < k; b++ {
				mask |= 1 << (b % 32)
				rules = append(rules, dtree.Rule[int]{
					Mask:    mask,
					Match:   0,
					Payload: len(rules),
				})
			}

			tree := dtree.New(rules)
			require.Equal(t, 1, tree.Nodes())
			require.Equal(t, 1, tree.Leaves())
			require.Equal(t, 1, tree.MaxDepth())
		})
	}
}
