package arm64

// Generator for addp — one constructor (Addp) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddpParams — parameters of addp.
type AddpParams struct {
	V3Params
}

func NewAddpParams(p V3Params) AddpParams {
	return AddpParams{V3Params: p}
}

func (p AddpParams) Instr() arm64.Instr {
	in, err := arm64.New().Addp(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddpParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Addp — an arbitrary addp.
func Addp(rnd *rand.Rand) ohsnap.Arbitrary[AddpParams] {
	base := newV3Gen(rnd, arrFull())
	return addpArb{base: base}
}

type addpArb struct {
	base v3Gen
}

func (a addpArb) Generate() iter.Seq[AddpParams] {
	return arbStream(func() AddpParams {
		return NewAddpParams(ohsnap.First(a.base.Generate()))
	})
}

func (a addpArb) Shrink(p AddpParams) iter.Seq[AddpParams] {
	var out []AddpParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewAddpParams(s))
	}

	return slices.Values(out)
}
