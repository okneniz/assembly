package alias

// Generator for the sbfiz alias — one generator, one type, one text form
// family: sbfiz rd, rn, #lsb, #width.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// SbfizParams — parameters of the sbfiz alias.
type SbfizParams struct {
	Rd, Rn arm64.Reg
	Lsb    uint32
	Width  uint32
}

func NewSbfizParams(rd arm64.Reg, rn arm64.Reg, lsb uint32, width uint32) SbfizParams {
	return SbfizParams{
		Rd:    rd,
		Rn:    rn,
		Lsb:   lsb,
		Width: width,
	}
}

func (p SbfizParams) String() string {
	return "sbfiz " + p.Rd.String() + ", " + p.Rn.String() +
		", #" + utoa(p.Lsb) + ", #" + utoa(p.Width)
}

func (p SbfizParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// sbfizGen — generator for sbfiz: same-width registers; lsb < size,
// 1 <= width, lsb+width <= size.
type sbfizGen struct {
	rnd *rand.Rand
}

func newSbfizGen(rnd *rand.Rand) sbfizGen {
	return sbfizGen{rnd: rnd}
}

// Sbfiz — an arbitrary sbfiz.
func Sbfiz(rnd *rand.Rand) ohsnap.Arbitrary[SbfizParams] {
	return newSbfizGen(rnd)
}

func (g sbfizGen) Generate() iter.Seq[SbfizParams] {
	return stream(func() SbfizParams {
		is64 := g.rnd.IntN(2) == 1
		size := uint32(32)
		if is64 {
			size = 64
		}

		lsb := g.rnd.IntN(int(size))
		return NewSbfizParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			uint32(lsb),
			uint32(g.rnd.IntN(int(size)-lsb))+1,
		)
	})
}

func (g sbfizGen) Shrink(p SbfizParams) iter.Seq[SbfizParams] {
	var out []SbfizParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewSbfizParams(r, p.Rn, p.Lsb, p.Width))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewSbfizParams(p.Rd, r, p.Lsb, p.Width))
	}

	for _, v := range uhalved(p.Lsb) {
		out = append(out, NewSbfizParams(p.Rd, p.Rn, v, p.Width))
	}

	for _, v := range uhalved(p.Width) {
		if v == 0 {
			continue // the width keeps >= 1
		}

		out = append(out, NewSbfizParams(p.Rd, p.Rn, p.Lsb, v))
	}

	return slices.Values(out)
}
