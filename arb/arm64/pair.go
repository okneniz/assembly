package arm64

// The register-pair family core: ldp/stp rt, rt2, [rn, #imm7<<scale]
// — both registers of one width (the scale follows it); ldpsw is the
// x-only twin at scale 2.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
)

// PairParams — the parameters of the pair forms.
type PairParams struct {
	Rt, Rt2, Rn arm64.Reg
	Off         arm64.Off
}

func NewPairParams(rt arm64.Reg, rt2 arm64.Reg, rn arm64.Reg, off arm64.Off) PairParams {
	return PairParams{
		Rt:  rt,
		Rt2: rt2,
		Rn:  rn,
		Off: off,
	}
}

// pairGen — the pair generator: registers of one width (the width mode
// comes from the family), the offset aligned to the pair scale.
type pairGen struct {
	rnd   *rand.Rand
	rt64  *bool
	scale *uint32 // pinned scale (Ldpsw: 2); nil: by the register width
}

func newPairGen(rnd *rand.Rand, rt64 *bool, scale *uint32) pairGen {
	return pairGen{
		rnd:   rnd,
		rt64:  rt64,
		scale: scale,
	}
}

// pair — the shared core of Ldp/Stp (either width, scale by width).
func pair(rnd *rand.Rand) pairGen {
	return newPairGen(rnd, nil, nil)
}

// pairX — the x-only core at the ldpsw scale.
func pairX(rnd *rand.Rand) pairGen {
	x := true
	s := uint32(2)
	return newPairGen(rnd, &x, &s)
}

func (g pairGen) Generate() iter.Seq[PairParams] {
	return arbStream(func() PairParams {
		is64 := g.is64()
		scale := g.curScale(is64)

		return NewPairParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, true, true, false),
			ohsnap.First(PairOff(g.rnd, scale).Generate()),
		)
	})
}

func (g pairGen) Shrink(p PairParams) iter.Seq[PairParams] {
	scale := g.curScale(p.Rt.Is64())

	regs := regShrunk(p.Rt)
	out := make([]PairParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewPairParams(r, p.Rt2, p.Rn, p.Off))
	}

	for _, r := range regShrunk(p.Rt2) {
		out = append(out, NewPairParams(p.Rt, r, p.Rn, p.Off))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewPairParams(p.Rt, p.Rt2, r, p.Off))
	}

	for _, o := range slices.Collect(PairOff(g.rnd, scale).Shrink(p.Off)) {
		out = append(out, NewPairParams(p.Rt, p.Rt2, p.Rn, o))
	}

	return slices.Values(out)
}

func (g pairGen) curScale(is64 bool) uint32 {
	if g.scale != nil {
		return *g.scale
	}

	if is64 {
		return 3
	}

	return 2
}

func (g pairGen) is64() bool {
	if g.rt64 == nil {
		return g.rnd.IntN(2) == 1
	}

	return *g.rt64
}
