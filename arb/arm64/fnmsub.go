package arm64

// Generator for fnmsub — one constructor (Fnmsub) over the shared
// four-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FnmsubParams — parameters of fnmsub.
type FnmsubParams struct {
	F4Params
}

func NewFnmsubParams(p F4Params) FnmsubParams {
	return FnmsubParams{F4Params: p}
}

func (p FnmsubParams) Instr() arm64.Instr {
	in, err := arm64.New().Fnmsub(p.Rd, p.Rn, p.Rm, p.Ra)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FnmsubParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fnmsub — an arbitrary fnmsub.
func Fnmsub(rnd *rand.Rand) ohsnap.Arbitrary[FnmsubParams] {
	base := f4(rnd)
	return fnmsubArb{base: base}
}

type fnmsubArb struct {
	base f4Gen
}

func (a fnmsubArb) Generate() iter.Seq[FnmsubParams] {
	return arbStream(func() FnmsubParams {
		return NewFnmsubParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fnmsubArb) Shrink(p FnmsubParams) iter.Seq[FnmsubParams] {
	shrinks := slices.Collect(a.base.Shrink(p.F4Params))
	out := make([]FnmsubParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewFnmsubParams(s))
	}

	return slices.Values(out)
}
