package dtree_test

import (
	mrnd "math/rand/v2"
	"testing"
	"time"

	"github.com/okneniz/assembly/dtree"
	"github.com/okneniz/assembly/dtree/arb"
)

// benchSink protects the traversals from being eliminated by the compiler.
var benchSink int

// benchField is the common defined fields of the "registry": sf[31],
// op[28:25], opc[23:21].
const benchField = 0x9E00_0000

// benchRegistry normalizes a live Generate rule to the registry structure:
// the sf/op/opc fields are additionally defined with random values. Without
// this, random broad masks cover each bit about half the time - duplicates
// are always above the 1/4 threshold and the tree degenerates into a leaf
// list: the benchmark would measure a linear scan, not a descent.
func benchRegistry(rnd *mrnd.Rand, r dtree.Rule[int]) dtree.Rule[int] {
	match := r.Match &^ benchField
	match |= (rnd.Uint32() & 0x1) << 31 // sf
	match |= (rnd.Uint32() & 0xF) << 25 // op
	match |= (rnd.Uint32() & 0x7) << 21 // opc

	return dtree.Rule[int]{
		Mask:    r.Mask | benchField,
		Match:   match,
		Payload: r.Payload,
	}
}

// benchCases is deterministic material from arb.LookupCase.Generate():
// cases (a rule set of <=40 + a word), a large "registry", and a word pool.
// The registry is live rules (without catch-alls and dead ones), normalized
// to the common fields; the words are half random (miss path) and half
// generated within the class of a random rule of the registry (hit path).
func benchCases(b *testing.B) (cases []arb.Case, big []dtree.Rule[int], words []uint32) {
	b.Helper()

	rnd := propSeed(b)
	gen := arb.LookupCase(rnd)

	for range 1000 {
		c := gen.Generate()
		cases = append(cases, c)

		for _, r := range c.Rules {
			if r.Mask != 0 && r.Match&^r.Mask == 0 {
				big = append(big, benchRegistry(rnd, r))
			}
		}

		if len(big) > 0 && rnd.IntN(2) == 0 {
			r := big[rnd.IntN(len(big))]
			words = append(words, r.Match|(rnd.Uint32()&^r.Mask))

			continue
		}

		words = append(words, rnd.Uint32())
	}

	return cases, big, words
}

// BenchmarkNew is tree construction from a random set of <=40 rules.
func BenchmarkNew(b *testing.B) {
	cases, _, _ := benchCases(b)

	b.ReportAllocs()
	b.ResetTimer()

	start := time.Now()
	var tree *dtree.Tree[int]
	for i := range b.N {
		tree = dtree.New(cases[i%len(cases)].Rules)
	}

	_ = tree

	b.ReportMetric(float64(b.N)/time.Since(start).Seconds(), "tree/s")
}

// BenchmarkNewLarge is construction from a concatenated set (~a thousand
// rules; leaf lists from catch-alls are part of the material).
func BenchmarkNewLarge(b *testing.B) {
	_, big, _ := benchCases(b)

	sample := dtree.New(big)
	b.Logf("rules: %d, nodes: %d, leaves: %d, depth: %d",
		len(big), sample.Nodes(), sample.Leaves(), sample.MaxDepth())

	b.ReportAllocs()
	b.ResetTimer()

	start := time.Now()
	var tree *dtree.Tree[int]
	for range b.N {
		tree = dtree.New(big)
	}

	_ = tree

	b.ReportMetric(float64(len(big))*float64(b.N)/time.Since(start).Seconds(), "rule/s")
}

// BenchmarkLookup is a descent of a large tree (random words).
func BenchmarkLookup(b *testing.B) {
	_, big, words := benchCases(b)
	tree := dtree.New(big)

	b.ReportAllocs()
	b.ResetTimer()

	start := time.Now()
	hits := 0
	for i := range b.N {
		if _, ok := tree.Lookup(words[i%len(words)]); ok {
			hits++
		}
	}

	_ = hits

	b.ReportMetric(float64(b.N)/time.Since(start).Seconds(), "lookup/s")
}

// BenchmarkNodes is a traversal of a large tree (the number of nodes).
func BenchmarkNodes(b *testing.B) {
	_, big, _ := benchCases(b)
	tree := dtree.New(big)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchSink = tree.Nodes()
	}
}

// BenchmarkLeaves is a traversal of a large tree (the number of leaves).
func BenchmarkLeaves(b *testing.B) {
	_, big, _ := benchCases(b)
	tree := dtree.New(big)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchSink = tree.Leaves()
	}
}

// BenchmarkMaxDepth is a traversal of a large tree (depth).
func BenchmarkMaxDepth(b *testing.B) {
	_, big, _ := benchCases(b)
	tree := dtree.New(big)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchSink = tree.MaxDepth()
	}
}
