package arm64

// The pc-relative offset generator: uniform in the checked byte range,
// always 4-byte aligned (the encode grid of every branch); the shrink
// halves toward 0 keeping the sign and the alignment — the range edges
// live at the boundaries.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
)

// brOffArb — a signed 4-byte-aligned offset in [-lim, lim) bytes.
type brOffArb struct {
	rnd *rand.Rand
	lim int64 // the exclusive byte bound (1<<20 imm19/imm21, 1<<26 imm26, 1<<15 imm14)
}

func newBrOffArb(rnd *rand.Rand, lim int64) brOffArb {
	return brOffArb{
		rnd: rnd,
		lim: lim,
	}
}

// BrOff — an arbitrary branch offset in the ±lim byte range.
func BrOff(rnd *rand.Rand, lim int64) ohsnap.Arbitrary[int64] {
	return newBrOffArb(rnd, lim)
}

func (g brOffArb) Generate() iter.Seq[int64] {
	n := g.lim / 4
	return arb.Stream(func() int64 {
		return (g.rnd.Int64N(2*n) - n) * 4
	})
}

func (g brOffArb) Shrink(off int64) iter.Seq[int64] {
	if off == 0 {
		return ohsnap.Empty[int64]()
	}

	var out []int64
	// the range edges first (decoder bugs live there), then halving
	// toward zero keeping the sign; the alignment survives every step.
	for _, v := range slices.Collect(shrink.Boundaries[int64](-(g.lim - 4), g.lim-4)(off)) {
		out = append(out, v)
	}

	for d := off / 2 / 4 * 4; ; d = d / 2 / 4 * 4 {
		out = append(out, d)
		if d == 0 {
			break
		}
	}

	return slices.Values(out)
}
