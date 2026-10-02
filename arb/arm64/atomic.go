package arm64

// The atomic load/store family core: ldar/ldaxr/stlr (rt x/w, rn x/sp)
// and the b halves (rt w) — the exclusive pair (stlxr/stxrb) adds the w
// status register.

import (
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/assembly/arch/arm64"
)

// ExclParams — the parameters of stlxr/stxrb rs, rt, [rn].
type ExclParams struct {
	Rs, Rt, Rn arm64.Reg
}

func NewExclParams(rs arm64.Reg, rt arm64.Reg, rn arm64.Reg) ExclParams {
	return ExclParams{
		Rs: rs,
		Rt: rt,
		Rn: rn,
	}
}

// atomGen — the atomic (rt, rn) generator; the width mode comes from
// the family.
type atomGen struct {
	rnd  *rand.Rand
	rt64 *bool
}

func newAtomGen(rnd *rand.Rand, rt64 *bool) atomGen {
	return atomGen{
		rnd:  rnd,
		rt64: rt64,
	}
}

// atom — the shared core of Ldar/Ldaxr/Stlr (either width).
func atom(rnd *rand.Rand) atomGen {
	return newAtomGen(rnd, nil)
}

// atomW — the w-only core (the b forms).
func atomW(rnd *rand.Rand) atomGen {
	w := false
	return newAtomGen(rnd, &w)
}

func (g atomGen) is64() bool {
	if g.rt64 == nil {
		return g.rnd.IntN(2) == 1
	}

	return *g.rt64
}

func (g atomGen) Generate() iter.Seq[LsParams] {
	return arbStream(func() LsParams {
		return NewLsParams(
			genReg(g.rnd, g.is64(), false, true),
			genReg(g.rnd, true, true, false),
			0,
		)
	})
}

func (g atomGen) Shrink(p LsParams) iter.Seq[LsParams] {
	var out []LsParams
	for _, r := range regShrunk(p.Rt) {
		out = append(out, NewLsParams(r, p.Rn, 0))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewLsParams(p.Rt, r, 0))
	}

	return slices.Values(out)
}

// exclGen — the exclusive-pair generator: rs is a w status register,
// rt of the family's width, rn x/sp.
type exclGen struct {
	rnd  *rand.Rand
	rt64 *bool
}

func newExclGen(rnd *rand.Rand, rt64 *bool) exclGen {
	return exclGen{
		rnd:  rnd,
		rt64: rt64,
	}
}

// excl — the shared core of Stlxr/Stxrb (either width rt).
func excl(rnd *rand.Rand) exclGen {
	return newExclGen(rnd, nil)
}

// exclW — the w-rt core (Stlxrb).
func exclW(rnd *rand.Rand) exclGen {
	w := false
	return newExclGen(rnd, &w)
}

func (g exclGen) is64() bool {
	if g.rt64 == nil {
		return g.rnd.IntN(2) == 1
	}

	return *g.rt64
}

func (g exclGen) Generate() iter.Seq[ExclParams] {
	return arbStream(func() ExclParams {
		return NewExclParams(
			genReg(g.rnd, false, false, true),
			genReg(g.rnd, g.is64(), false, true),
			genReg(g.rnd, true, true, false),
		)
	})
}

func (g exclGen) Shrink(p ExclParams) iter.Seq[ExclParams] {
	var out []ExclParams
	for _, r := range regShrunk(p.Rs) {
		out = append(out, NewExclParams(r, p.Rt, p.Rn))
	}

	for _, r := range regShrunk(p.Rt) {
		out = append(out, NewExclParams(p.Rs, r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewExclParams(p.Rs, p.Rt, r))
	}

	return slices.Values(out)
}
