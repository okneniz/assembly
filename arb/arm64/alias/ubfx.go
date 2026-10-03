package alias

// Generator for the ubfx alias — one generator, one type, one text form
// family: ubfx rd, rn, #lsb, #width.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// UbfxParams — parameters of the ubfx alias.
type UbfxParams struct {
	Rd, Rn arm64.Reg
	Lsb    uint32
	Width  uint32
}

func NewUbfxParams(rd arm64.Reg, rn arm64.Reg, lsb uint32, width uint32) UbfxParams {
	return UbfxParams{
		Rd:    rd,
		Rn:    rn,
		Lsb:   lsb,
		Width: width,
	}
}

func (p UbfxParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p UbfxParams) String() string {
	return "ubfx " + p.Rd.String() + ", " + p.Rn.String() +
		", #" + utoa(p.Lsb) + ", #" + utoa(p.Width)
}

// ubfxGen — generator for ubfx: same-width registers; lsb < size,
// 1 <= width, lsb+width <= size.
type ubfxGen struct {
	rnd *rand.Rand
}

// Ubfx — an arbitrary ubfx.
func Ubfx(rnd *rand.Rand) ohsnap.Arbitrary[UbfxParams] {
	return newUbfxGen(rnd)
}

func newUbfxGen(rnd *rand.Rand) ubfxGen {
	return ubfxGen{rnd: rnd}
}

func (g ubfxGen) Generate() iter.Seq[UbfxParams] {
	return stream(func() UbfxParams {
		is64 := g.rnd.IntN(2) == 1
		size := uint32(32)
		if is64 {
			size = 64
		}

		lsb := g.rnd.IntN(int(size))
		return NewUbfxParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			uint32(lsb),
			uint32(g.rnd.IntN(int(size)-lsb))+1,
		)
	})
}

func (g ubfxGen) Shrink(p UbfxParams) iter.Seq[UbfxParams] {
	var out []UbfxParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewUbfxParams(r, p.Rn, p.Lsb, p.Width))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewUbfxParams(p.Rd, r, p.Lsb, p.Width))
	}

	for _, v := range uhalved(p.Lsb) {
		out = append(out, NewUbfxParams(p.Rd, p.Rn, v, p.Width))
	}

	for _, v := range uhalved(p.Width) {
		if v == 0 {
			continue // the width keeps >= 1
		}

		out = append(out, NewUbfxParams(p.Rd, p.Rn, p.Lsb, v))
	}

	return slices.Values(out)
}
