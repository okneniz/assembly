package arm64

// Generator for the ands logical immediate — one constructor (AndsImm)
// over the shared bitmask core (the structural esize/len/rot axes).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AndsImmParams — parameters of ands rd, rn, #bitmask.
type AndsImmParams struct {
	BitmaskParams
}

func NewAndsImmParams(p BitmaskParams) AndsImmParams {
	return AndsImmParams{BitmaskParams: p}
}

func (p AndsImmParams) Instr() arm64.Instr {
	in, err := arm64.New().AndsImm(p.Rd, p.Rn, p.Value())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AndsImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AndsImm — an arbitrary ands (logical immediate).
func AndsImm(rnd *rand.Rand) ohsnap.Arbitrary[AndsImmParams] {
	base := bitmask(rnd)
	return andsImmArb{base: base}
}

type andsImmArb struct {
	base bitmaskGen
}

func (a andsImmArb) Generate() iter.Seq[AndsImmParams] {
	return arbStream(func() AndsImmParams {
		return NewAndsImmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a andsImmArb) Shrink(p AndsImmParams) iter.Seq[AndsImmParams] {
	var out []AndsImmParams
	for _, s := range slices.Collect(a.base.Shrink(p.BitmaskParams)) {
		out = append(out, NewAndsImmParams(s))
	}

	return slices.Values(out)
}
