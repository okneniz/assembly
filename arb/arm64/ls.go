package arm64

// The load/store (imm) family core: one parameter shape {rt, rn, off};
// the register-width class of rt and the offset scale/format are the
// family's own business — the thin per-instruction files bind them.
// rt: x/w (the 31st reads as zr), rn: x (the 31st reads as sp).

import (
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/assembly/arch/arm64"
	ohsnap "github.com/okneniz/oh-snap"
)

// LsParams — the parameters of the imm-addressed load/store forms.
type LsParams struct {
	Rt, Rn arm64.Reg
	Off    arm64.Off
}

func NewLsParams(rt arm64.Reg, rn arm64.Reg, off arm64.Off) LsParams {
	return LsParams{
		Rt:  rt,
		Rn:  rn,
		Off: off,
	}
}

// lsGen — the shared generator: rt width mode (nil: either) and the
// offset generator come from the family.
type lsGen struct {
	rnd  *rand.Rand
	rt64 *bool // nil: either width; else pinned
	off  ohsnapOff
}

// ohsnapOff — the offset arbitrary interface (Generate/Shrink).
type ohsnapOff interface {
	Generate() iter.Seq[arm64.Off]
	Shrink(arm64.Off) iter.Seq[arm64.Off]
}

func newLsGen(rnd *rand.Rand, rt64 *bool, off ohsnapOff) lsGen {
	return lsGen{
		rnd:  rnd,
		rt64: rt64,
		off:  off,
	}
}

func (g lsGen) is64() bool {
	if g.rt64 == nil {
		return g.rnd.IntN(2) == 1
	}

	return *g.rt64
}

func (g lsGen) Generate() iter.Seq[LsParams] {
	return arbStream(func() LsParams {
		is64 := g.is64()
		return NewLsParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, true, true, false),
			ohsnap.First(g.off.Generate()),
		)
	})
}

func (g lsGen) Shrink(p LsParams) iter.Seq[LsParams] {
	var out []LsParams
	for _, r := range regShrunk(p.Rt) {
		out = append(out, NewLsParams(r, p.Rn, p.Off))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewLsParams(p.Rt, r, p.Off))
	}

	for _, o := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewLsParams(p.Rt, p.Rn, o))
	}

	return slices.Values(out)
}
