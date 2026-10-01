package arm64

// Generator for the eor logical immediate — one constructor (EorImm)
// over the shared bitmask core (the structural esize/len/rot axes).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// EorImmParams — parameters of eor rd, rn, #bitmask.
type EorImmParams struct {
	BitmaskParams
}

func NewEorImmParams(p BitmaskParams) EorImmParams {
	return EorImmParams{BitmaskParams: p}
}

func (p EorImmParams) Instr() arm64.Instr {
	in, err := arm64.New().EorImm(p.Rd, p.Rn, p.Value())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p EorImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// EorImm — an arbitrary eor (logical immediate).
func EorImm(rnd *rand.Rand) ohsnap.Arbitrary[EorImmParams] {
	base := bitmask(rnd)
	return eorImmArb{base: base}
}

type eorImmArb struct {
	base bitmaskGen
}

func (a eorImmArb) Generate() iter.Seq[EorImmParams] {
	return arbStream(func() EorImmParams {
		return NewEorImmParams(ohsnap.First(a.base.Generate()))
	})
}

func (a eorImmArb) Shrink(p EorImmParams) iter.Seq[EorImmParams] {
	var out []EorImmParams
	for _, s := range slices.Collect(a.base.Shrink(p.BitmaskParams)) {
		out = append(out, NewEorImmParams(s))
	}

	return slices.Values(out)
}
