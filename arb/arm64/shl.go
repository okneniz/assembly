package arm64

// Generator for shl — one constructor (Shl) over the shared
// element-bound shift core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ShlParams — parameters of shl.
type ShlParams struct {
	VShiftParams
}

func NewShlParams(p VShiftParams) ShlParams {
	return ShlParams{VShiftParams: p}
}

func (p ShlParams) Instr() arm64.Instr {
	in, err := arm64.New().Shl(p.Rd, p.Rn, p.Arr, p.Shift)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p ShlParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Shl — an arbitrary shl.
func Shl(rnd *rand.Rand) ohsnap.Arbitrary[ShlParams] {
	base := newVShiftGen(rnd, arrFull())
	return shlArb{base: base}
}

type shlArb struct {
	base vShiftGen
}

func (a shlArb) Generate() iter.Seq[ShlParams] {
	return arbStream(func() ShlParams {
		return NewShlParams(ohsnap.First(a.base.Generate()))
	})
}

func (a shlArb) Shrink(p ShlParams) iter.Seq[ShlParams] {
	var out []ShlParams
	for _, s := range slices.Collect(a.base.Shrink(p.VShiftParams)) {
		out = append(out, NewShlParams(s))
	}

	return slices.Values(out)
}
