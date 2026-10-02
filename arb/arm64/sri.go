package arm64

// Generator for sri — one constructor (Sri) over the shared
// element-bound shift core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SriParams — parameters of sri.
type SriParams struct {
	VShiftParams
}

func NewSriParams(p VShiftParams) SriParams {
	return SriParams{VShiftParams: p}
}

func (p SriParams) Instr() arm64.Instr {
	in, err := arm64.New().Sri(p.Rd, p.Rn, p.Arr, p.Shift)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SriParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Sri — an arbitrary sri.
func Sri(rnd *rand.Rand) ohsnap.Arbitrary[SriParams] {
	base := newVShiftGen1(rnd, arrFull())
	return sriArb{base: base}
}

type sriArb struct {
	base vShiftGen
}

func (a sriArb) Generate() iter.Seq[SriParams] {
	return arbStream(func() SriParams {
		return NewSriParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sriArb) Shrink(p SriParams) iter.Seq[SriParams] {
	var out []SriParams
	for _, s := range slices.Collect(a.base.Shrink(p.VShiftParams)) {
		out = append(out, NewSriParams(s))
	}

	return slices.Values(out)
}
