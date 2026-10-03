package arm64

// Generator for adds (extended register) — one constructor (AddsExt) over
// the shared ext core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddsExtParams — parameters of adds rd, rn, rm, ext #imm3.
type AddsExtParams struct {
	ExtParams
}

func NewAddsExtParams(p ExtParams) AddsExtParams {
	return AddsExtParams{ExtParams: p}
}

func (p AddsExtParams) Instr() arm64.Instr {
	in, err := arm64.New().AddsExt(p.Rd, p.Rn, p.Rm, p.Ext, p.Imm3)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddsExtParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AddsExt — an arbitrary adds (extended register).
func AddsExt(rnd *rand.Rand) ohsnap.Arbitrary[AddsExtParams] {
	base := extS(rnd)
	return addsExtArb{base: base}
}

type addsExtArb struct {
	base extGen
}

func (a addsExtArb) Generate() iter.Seq[AddsExtParams] {
	return arbStream(func() AddsExtParams {
		return NewAddsExtParams(ohsnap.First(a.base.Generate()))
	})
}

func (a addsExtArb) Shrink(p AddsExtParams) iter.Seq[AddsExtParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ExtParams))
	out := make([]AddsExtParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewAddsExtParams(s))
	}

	return slices.Values(out)
}
