package arm64

// Generator for fdiv — one constructor (Fdiv) over the shared
// three-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FdivParams — parameters of fdiv.
type FdivParams struct {
	F3Params
}

func NewFdivParams(p F3Params) FdivParams {
	return FdivParams{F3Params: p}
}

func (p FdivParams) Instr() arm64.Instr {
	in, err := arm64.New().Fdiv(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FdivParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fdiv — an arbitrary fdiv.
func Fdiv(rnd *rand.Rand) ohsnap.Arbitrary[FdivParams] {
	base := f3(rnd)
	return fdivArb{base: base}
}

type fdivArb struct {
	base f3Gen
}

func (a fdivArb) Generate() iter.Seq[FdivParams] {
	return arbStream(func() FdivParams {
		return NewFdivParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fdivArb) Shrink(p FdivParams) iter.Seq[FdivParams] {
	var out []FdivParams
	for _, s := range slices.Collect(a.base.Shrink(p.F3Params)) {
		out = append(out, NewFdivParams(s))
	}

	return slices.Values(out)
}
