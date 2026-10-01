package arm64

// Generator for the orr logical immediate — one constructor (OrrImm)
// over the shared bitmask core (the structural esize/len/rot axes).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// OrrImmParams — parameters of orr rd, rn, #bitmask.
type OrrImmParams struct {
	BitmaskParams
}

func NewOrrImmParams(p BitmaskParams) OrrImmParams {
	return OrrImmParams{BitmaskParams: p}
}

func (p OrrImmParams) Instr() arm64.Instr {
	in, err := arm64.New().OrrImm(p.Rd, p.Rn, p.Value())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p OrrImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// OrrImm — an arbitrary orr (logical immediate).
func OrrImm(rnd *rand.Rand) ohsnap.Arbitrary[OrrImmParams] {
	base := bitmask(rnd)
	return orrImmArb{base: base}
}

type orrImmArb struct {
	base bitmaskGen
}

func (a orrImmArb) Generate() iter.Seq[OrrImmParams] {
	return arbStream(func() OrrImmParams {
		return NewOrrImmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a orrImmArb) Shrink(p OrrImmParams) iter.Seq[OrrImmParams] {
	var out []OrrImmParams
	for _, s := range slices.Collect(a.base.Shrink(p.BitmaskParams)) {
		out = append(out, NewOrrImmParams(s))
	}

	return slices.Values(out)
}
