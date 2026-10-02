package arm64

// Generator for strh — one constructor (Strh) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StrhParams — parameters of the strh form.
type StrhParams struct {
	LsParams
}

func NewStrhParams(p LsParams) StrhParams {
	return StrhParams{LsParams: p}
}

func (p StrhParams) Instr() arm64.Instr {
	in, err := arm64.New().Strh(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StrhParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Strh — an arbitrary strh.
func Strh(rnd *rand.Rand) ohsnap.Arbitrary[StrhParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, ScaledOff(rnd, 1))
	return strhArb{base: base}
}

type strhArb struct {
	base lsGen
}

func (a strhArb) Generate() iter.Seq[StrhParams] {
	return arbStream(func() StrhParams {
		return NewStrhParams(ohsnap.First(a.base.Generate()))
	})
}

func (a strhArb) Shrink(p StrhParams) iter.Seq[StrhParams] {
	var out []StrhParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewStrhParams(s))
	}

	return slices.Values(out)
}
