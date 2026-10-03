package alias

// Generator for the ubfiz alias — one generator, one type, one text form
// family: ubfiz rd, rn, #lsb, #width.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// UbfizParams — parameters of the ubfiz alias.
type UbfizParams struct {
	Rd, Rn arm64.Reg
	Lsb    uint32
	Width  uint32
}

func NewUbfizParams(rd arm64.Reg, rn arm64.Reg, lsb uint32, width uint32) UbfizParams {
	return UbfizParams{
		Rd:    rd,
		Rn:    rn,
		Lsb:   lsb,
		Width: width,
	}
}

func (p UbfizParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p UbfizParams) String() string {
	return "ubfiz " + p.Rd.String() + ", " + p.Rn.String() +
		", #" + utoa(p.Lsb) + ", #" + utoa(p.Width)
}

// ubfizGen — generator for ubfiz: same-width registers; lsb < size,
// 1 <= width, lsb+width <= size.
type ubfizGen struct {
	rnd *rand.Rand
}

// Ubfiz — an arbitrary ubfiz.
func Ubfiz(rnd *rand.Rand) ohsnap.Arbitrary[UbfizParams] {
	return newUbfizGen(rnd)
}

func newUbfizGen(rnd *rand.Rand) ubfizGen {
	return ubfizGen{rnd: rnd}
}

func (g ubfizGen) Generate() iter.Seq[UbfizParams] {
	return stream(func() UbfizParams {
		is64 := g.rnd.IntN(2) == 1
		size := uint32(32)
		if is64 {
			size = 64
		}

		lsb := g.rnd.IntN(int(size))
		return NewUbfizParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			uint32(lsb),
			uint32(g.rnd.IntN(int(size)-lsb))+1,
		)
	})
}

func (g ubfizGen) Shrink(p UbfizParams) iter.Seq[UbfizParams] {
	var out []UbfizParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewUbfizParams(r, p.Rn, p.Lsb, p.Width))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewUbfizParams(p.Rd, r, p.Lsb, p.Width))
	}

	for _, v := range uhalved(p.Lsb) {
		out = append(out, NewUbfizParams(p.Rd, p.Rn, v, p.Width))
	}

	for _, v := range uhalved(p.Width) {
		if v == 0 {
			continue // the width keeps >= 1
		}

		out = append(out, NewUbfizParams(p.Rd, p.Rn, p.Lsb, v))
	}

	return slices.Values(out)
}
