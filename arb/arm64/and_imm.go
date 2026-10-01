package arm64

// Generator for the and logical immediate — one constructor (AndImm)
// over the shared bitmask core (the structural esize/len/rot axes).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AndImmParams — parameters of and rd, rn, #bitmask.
type AndImmParams struct {
	BitmaskParams
}

func NewAndImmParams(p BitmaskParams) AndImmParams {
	return AndImmParams{BitmaskParams: p}
}

func (p AndImmParams) Instr() arm64.Instr {
	in, err := arm64.New().AndImm(p.Rd, p.Rn, p.Value())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AndImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AndImm — an arbitrary and (logical immediate).
func AndImm(rnd *rand.Rand) ohsnap.Arbitrary[AndImmParams] {
	base := bitmask(rnd)
	return andImmArb{base: base}
}

type andImmArb struct {
	base bitmaskGen
}

func (a andImmArb) Generate() iter.Seq[AndImmParams] {
	return arbStream(func() AndImmParams {
		return NewAndImmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a andImmArb) Shrink(p AndImmParams) iter.Seq[AndImmParams] {
	var out []AndImmParams
	for _, s := range slices.Collect(a.base.Shrink(p.BitmaskParams)) {
		out = append(out, NewAndImmParams(s))
	}

	return slices.Values(out)
}
