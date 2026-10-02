package arm64

// Generator for add — one constructor (Add) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddVParams — parameters of add.
type AddVParams struct {
	V3Params
}

func NewAddVParams(p V3Params) AddVParams {
	return AddVParams{V3Params: p}
}

func (p AddVParams) Instr() arm64.Instr {
	in, err := arm64.New().Add(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddVParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Add — an arbitrary add.
func Add(rnd *rand.Rand) ohsnap.Arbitrary[AddVParams] {
	base := newV3Gen(rnd, arrFull())
	return addArb{base: base}
}

type addArb struct {
	base v3Gen
}

func (a addArb) Generate() iter.Seq[AddVParams] {
	return arbStream(func() AddVParams {
		return NewAddVParams(ohsnap.First(a.base.Generate()))
	})
}

func (a addArb) Shrink(p AddVParams) iter.Seq[AddVParams] {
	var out []AddVParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewAddVParams(s))
	}

	return slices.Values(out)
}
