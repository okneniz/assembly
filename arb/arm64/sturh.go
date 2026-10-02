package arm64

// Generator for sturh — one constructor (Sturh) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SturhParams — parameters of the sturh form.
type SturhParams struct {
	LsParams
}

func NewSturhParams(p LsParams) SturhParams {
	return SturhParams{LsParams: p}
}

func (p SturhParams) Instr() arm64.Instr {
	in, err := arm64.New().Sturh(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SturhParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Sturh — an arbitrary sturh.
func Sturh(rnd *rand.Rand) ohsnap.Arbitrary[SturhParams] {
	rtW := false
	base := newLsGen(rnd, &rtW, UnscaledOff(rnd))
	return sturhArb{base: base}
}

type sturhArb struct {
	base lsGen
}

func (a sturhArb) Generate() iter.Seq[SturhParams] {
	return arbStream(func() SturhParams {
		return NewSturhParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sturhArb) Shrink(p SturhParams) iter.Seq[SturhParams] {
	var out []SturhParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewSturhParams(s))
	}

	return slices.Values(out)
}
