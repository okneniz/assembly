package alias

// Generator for the sbfx alias — one generator, one type, one text form
// family: sbfx rd, rn, #lsb, #width.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// SbfxParams — parameters of the sbfx alias.
type SbfxParams struct {
	Rd, Rn arm64.Reg
	Lsb    uint32
	Width  uint32
}

func NewSbfxParams(rd arm64.Reg, rn arm64.Reg, lsb uint32, width uint32) SbfxParams {
	return SbfxParams{
		Rd:    rd,
		Rn:    rn,
		Lsb:   lsb,
		Width: width,
	}
}

func (p SbfxParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p SbfxParams) String() string {
	return "sbfx " + p.Rd.String() + ", " + p.Rn.String() +
		", #" + utoa(p.Lsb) + ", #" + utoa(p.Width)
}

// sbfxGen — generator for sbfx: same-width registers; lsb < size,
// 1 <= width, lsb+width <= size.
type sbfxGen struct {
	rnd *rand.Rand
}

// Sbfx — an arbitrary sbfx.
func Sbfx(rnd *rand.Rand) ohsnap.Arbitrary[SbfxParams] {
	return newSbfxGen(rnd)
}

func newSbfxGen(rnd *rand.Rand) sbfxGen {
	return sbfxGen{rnd: rnd}
}

func (g sbfxGen) Generate() iter.Seq[SbfxParams] {
	return stream(func() SbfxParams {
		is64 := g.rnd.IntN(2) == 1
		size := uint32(32)
		if is64 {
			size = 64
		}

		lsb := g.rnd.IntN(int(size))
		return NewSbfxParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			uint32(lsb),
			uint32(g.rnd.IntN(int(size)-lsb))+1,
		)
	})
}

func (g sbfxGen) Shrink(p SbfxParams) iter.Seq[SbfxParams] {
	var out []SbfxParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewSbfxParams(r, p.Rn, p.Lsb, p.Width))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewSbfxParams(p.Rd, r, p.Lsb, p.Width))
	}

	for _, v := range uhalved(p.Lsb) {
		out = append(out, NewSbfxParams(p.Rd, p.Rn, v, p.Width))
	}

	for _, v := range uhalved(p.Width) {
		if v == 0 {
			continue // the width keeps >= 1
		}

		out = append(out, NewSbfxParams(p.Rd, p.Rn, p.Lsb, v))
	}

	return slices.Values(out)
}
