package arm64

// The load/store offset generators: the scaled imm12 forms (0..0xfff
// shifted by the element scale, aligned) and the unscaled imm9 forms
// (-256..255, unaligned). The shrink keeps the alignment and halves
// toward 0 — the range edges live at the boundaries.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
)

// scaledOffArb — an aligned imm12 offset of the given element scale.
type scaledOffArb struct {
	rnd   *rand.Rand
	scale uint32
}

func newScaledOffArb(rnd *rand.Rand, scale uint32) scaledOffArb {
	return scaledOffArb{
		rnd:   rnd,
		scale: scale,
	}
}

// ScaledOff — an arbitrary aligned load/store offset (0..0xfff << scale).
func ScaledOff(rnd *rand.Rand, scale uint32) ohsnap.Arbitrary[arm64.Off] {
	return newScaledOffArb(rnd, scale)
}

func (a scaledOffArb) Generate() iter.Seq[arm64.Off] {
	return arb.Stream(func() arm64.Off {
		return arm64.Off(a.rnd.Int64N(0x1000) << a.scale)
	})
}

func (a scaledOffArb) Shrink(o arm64.Off) iter.Seq[arm64.Off] {
	return slices.Values(offShrunk(o, a.scale))
}

// unscaledOffArb — an imm9 offset (-256..255), no alignment.
type unscaledOffArb struct {
	rnd *rand.Rand
}

func newUnscaledOffArb(rnd *rand.Rand) unscaledOffArb {
	return unscaledOffArb{rnd: rnd}
}

// UnscaledOff — an arbitrary unscaled load/store offset (-256..255).
func UnscaledOff(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Off] {
	return newUnscaledOffArb(rnd)
}

func (a unscaledOffArb) Generate() iter.Seq[arm64.Off] {
	return arb.Stream(func() arm64.Off {
		return arm64.Off(a.rnd.Int64N(512) - 256)
	})
}

func (a unscaledOffArb) Shrink(o arm64.Off) iter.Seq[arm64.Off] {
	var out []arm64.Off
	for d := range shrink.Boundaries[int64](-256, 255)(int64(o)) {
		out = append(out, arm64.Off(d))
	}

	for d := range shrink.Halving[int64](0)(int64(o)) {
		out = append(out, arm64.Off(d))
	}

	return slices.Values(out)
}

// pairOffArb — a pair offset (imm7 shifted by the pair scale, aligned,
// -64..63 elements).
type pairOffArb struct {
	rnd   *rand.Rand
	scale uint32
}

func newPairOffArb(rnd *rand.Rand, scale uint32) pairOffArb {
	return pairOffArb{
		rnd:   rnd,
		scale: scale,
	}
}

// PairOff — an arbitrary pair load/store offset (-64..63 << scale).
func PairOff(rnd *rand.Rand, scale uint32) ohsnap.Arbitrary[arm64.Off] {
	return newPairOffArb(rnd, scale)
}

func (a pairOffArb) Generate() iter.Seq[arm64.Off] {
	return arb.Stream(func() arm64.Off {
		return arm64.Off((a.rnd.Int64N(128) - 64) << a.scale)
	})
}

func (a pairOffArb) Shrink(o arm64.Off) iter.Seq[arm64.Off] {
	var out []arm64.Off
	for d := range shrink.Boundaries[int64](-64<<a.scale, 63<<a.scale)(int64(o)) {
		out = append(out, arm64.Off(d))
	}

	for d := range shrink.Halving[int64](0)(int64(o)) {
		out = append(out, arm64.Off(d))
	}

	return slices.Values(out)
}
