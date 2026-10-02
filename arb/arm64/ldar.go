package arm64

// Generator for ldar — one constructor (Ldar) over the shared
// atomic core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdarParams — parameters of the ldar form.
type LdarParams struct {
	LsParams
}

func NewLdarParams(p LsParams) LdarParams {
	return LdarParams{LsParams: p}
}

func (p LdarParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldar(p.Rt, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdarParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldar — an arbitrary ldar.
func Ldar(rnd *rand.Rand) ohsnap.Arbitrary[LdarParams] {
	base := atom(rnd)
	return ldarArb{base: base}
}

type ldarArb struct {
	base atomGen
}

func (a ldarArb) Generate() iter.Seq[LdarParams] {
	return arbStream(func() LdarParams {
		return NewLdarParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldarArb) Shrink(p LdarParams) iter.Seq[LdarParams] {
	var out []LdarParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdarParams(s))
	}

	return slices.Values(out)
}
