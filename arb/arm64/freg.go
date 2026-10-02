package arm64

// The FP load/store family core: ldr/str st|dt, [xn, #off] — the FP
// register file has no named 31st, the offset scale follows the width
// (s: 2, d: 3).

import (
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/assembly/arch/arm64"
	ohsnap "github.com/okneniz/oh-snap"
)

// FParams — the parameters of the FP-addressed load/store forms.
type FParams struct {
	Rt  arm64.FReg
	Rn  arm64.Reg
	Off arm64.Off
}

func NewFParams(rt arm64.FReg, rn arm64.Reg, off arm64.Off) FParams {
	return FParams{
		Rt:  rt,
		Rn:  rn,
		Off: off,
	}
}

// fregGen — an FP register of either width: s0..s31 / d0..d31.
type fregGen struct {
	rnd *rand.Rand
}

func newFregGen(rnd *rand.Rand) fregGen {
	return fregGen{rnd: rnd}
}

// mustS/mustD — the checked FReg constructors: the input is bounded by
// the caller (0..31), so the error is unreachable.
func mustS(n int) arm64.FReg {
	r, err := arm64.S(n)
	if err != nil {
		return arm64.FReg{} // unreachable: n is always in 0..31
	}

	return r
}

func mustD(n int) arm64.FReg {
	r, err := arm64.D(n)
	if err != nil {
		return arm64.FReg{} // unreachable: n is always in 0..31
	}

	return r
}

func (g fregGen) Generate() iter.Seq[arm64.FReg] {
	return arbStream(func() arm64.FReg {
		n := g.rnd.IntN(32)
		if g.rnd.IntN(2) == 1 {
			return mustD(n)
		}

		return mustS(n)
	})
}

// Shrink — toward d0/s0, then the halved number of the same width.
func (g fregGen) Shrink(r arm64.FReg) iter.Seq[arm64.FReg] {
	var out []arm64.FReg
	zero := mustS(0)
	if r.Is64() {
		zero = mustD(0)
	}

	if r != zero {
		out = append(out, zero)
	}

	n := int(r.Num())
	if n > 1 {
		if r.Is64() {
			out = append(out, mustD(n/2))
		} else {
			out = append(out, mustS(n/2))
		}
	}

	return slices.Values(out)
}

// fGen — the FP load/store generator: an FP register, an x/sp base and
// the offset aligned to the register's element.
type fGen struct {
	rnd  *rand.Rand
	freg fregGen
}

func newFGen(rnd *rand.Rand) fGen {
	return fGen{
		rnd:  rnd,
		freg: newFregGen(rnd),
	}
}

// f — the shared core of LdrF/StrF.
func f(rnd *rand.Rand) fGen {
	return newFGen(rnd)
}

func (g fGen) Generate() iter.Seq[FParams] {
	return arbStream(func() FParams {
		rt := ohsnap.First(g.freg.Generate())
		scale := uint32(2)
		if rt.Is64() {
			scale = 3
		}

		return NewFParams(
			rt,
			genReg(g.rnd, true, true, false),
			ohsnap.First(ScaledOff(g.rnd, scale).Generate()),
		)
	})
}

func (g fGen) Shrink(p FParams) iter.Seq[FParams] {
	var out []FParams
	for _, r := range slices.Collect(g.freg.Shrink(p.Rt)) {
		out = append(out, NewFParams(r, p.Rn, p.Off))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewFParams(p.Rt, r, p.Off))
	}

	scale := uint32(2)
	if p.Rt.Is64() {
		scale = 3
	}

	for _, o := range slices.Collect(ScaledOff(g.rnd, scale).Shrink(p.Off)) {
		out = append(out, NewFParams(p.Rt, p.Rn, o))
	}

	return slices.Values(out)
}
