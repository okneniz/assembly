package alias

// Generator for the tst (immediate) alias — the logical immediate form
// over the structural bitmask core of arb/arm64 (the esize/len/rot axes
// generate the encodable value; the text renders it as clang does).

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// TstImmParams — parameters of the tst (immediate) alias: a register and
// the structural axes of the logical immediate.
type TstImmParams struct {
	Rn arm64.Reg // 31 reads as zr
	a64.BitmaskParams
}

func NewTstImmParams(rn arm64.Reg, p a64.BitmaskParams) TstImmParams {
	return TstImmParams{
		Rn:            rn,
		BitmaskParams: p,
	}
}

func (p TstImmParams) String() string {
	return fmt.Sprintf("tst %s, #0x%x", p.Rn, p.Value())
}

func (p TstImmParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// tstImmArb — the alias family over the shared axes: the register takes
// the width of the axes' registers (the whole space is one s/d-width
// file).
type tstImmArb struct {
	base ohsnap.Arbitrary[a64.BitmaskParams]
	rnd  *rand.Rand
}

// TstImm — an arbitrary tst (immediate).
func TstImm(rnd *rand.Rand) ohsnap.Arbitrary[TstImmParams] {
	return tstImmArb{
		base: a64.Bitmask(rnd),
		rnd:  rnd,
	}
}

func (a tstImmArb) Generate() iter.Seq[TstImmParams] {
	return stream(func() TstImmParams {
		p := ohsnap.First(a.base.Generate())
		return NewTstImmParams(
			a64.GenReg(a.rnd, p.Rd.Is64(), false, true),
			p,
		)
	})
}

func (a tstImmArb) Shrink(p TstImmParams) iter.Seq[TstImmParams] {
	var out []TstImmParams
	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewTstImmParams(r, p.BitmaskParams))
	}

	for _, s := range slices.Collect(a.base.Shrink(p.BitmaskParams)) {
		out = append(out, NewTstImmParams(p.Rn, s))
	}

	return slices.Values(out)
}
