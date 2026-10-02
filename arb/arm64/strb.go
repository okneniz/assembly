package arm64

// Generator for strb — one constructor (Strb) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StrbParams — parameters of the strb form.
type StrbParams struct {
	LsParams
}

func NewStrbParams(p LsParams) StrbParams {
	return StrbParams{LsParams: p}
}

func (p StrbParams) Instr() arm64.Instr {
	in, err := arm64.New().Strb(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StrbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Strb — an arbitrary strb.
func Strb(rnd *rand.Rand) ohsnap.Arbitrary[StrbParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, ScaledOff(rnd, 0))
	return strbArb{base: base}
}

type strbArb struct {
	base lsGen
}

func (a strbArb) Generate() iter.Seq[StrbParams] {
	return arbStream(func() StrbParams {
		return NewStrbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a strbArb) Shrink(p StrbParams) iter.Seq[StrbParams] {
	var out []StrbParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewStrbParams(s))
	}

	return slices.Values(out)
}
