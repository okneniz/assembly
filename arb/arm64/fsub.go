package arm64

// Generator for fsub — one constructor (Fsub) over the shared
// three-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FsubParams — parameters of fsub.
type FsubParams struct {
	F3Params
}

func NewFsubParams(p F3Params) FsubParams {
	return FsubParams{F3Params: p}
}

func (p FsubParams) Instr() arm64.Instr {
	in, err := arm64.New().Fsub(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FsubParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fsub — an arbitrary fsub.
func Fsub(rnd *rand.Rand) ohsnap.Arbitrary[FsubParams] {
	base := f3(rnd)
	return fsubArb{base: base}
}

type fsubArb struct {
	base f3Gen
}

func (a fsubArb) Generate() iter.Seq[FsubParams] {
	return arbStream(func() FsubParams {
		return NewFsubParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fsubArb) Shrink(p FsubParams) iter.Seq[FsubParams] {
	var out []FsubParams
	for _, s := range slices.Collect(a.base.Shrink(p.F3Params)) {
		out = append(out, NewFsubParams(s))
	}

	return slices.Values(out)
}
