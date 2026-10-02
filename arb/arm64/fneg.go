package arm64

// Generator for fneg — one constructor (Fneg) over the shared
// two-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FnegParams — parameters of fneg.
type FnegParams struct {
	F2Params
}

func NewFnegParams(p F2Params) FnegParams {
	return FnegParams{F2Params: p}
}

func (p FnegParams) Instr() arm64.Instr {
	in, err := arm64.New().Fneg(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p FnegParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fneg — an arbitrary fneg.
func Fneg(rnd *rand.Rand) ohsnap.Arbitrary[FnegParams] {
	base := f2(rnd)
	return fnegArb{base: base}
}

type fnegArb struct {
	base f2Gen
}

func (a fnegArb) Generate() iter.Seq[FnegParams] {
	return arbStream(func() FnegParams {
		return NewFnegParams(ohsnap.First(a.base.Generate()))
	})
}

func (a fnegArb) Shrink(p FnegParams) iter.Seq[FnegParams] {
	var out []FnegParams
	for _, s := range slices.Collect(a.base.Shrink(p.F2Params)) {
		out = append(out, NewFnegParams(s))
	}

	return slices.Values(out)
}
