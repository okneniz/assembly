package arm64

// Generator for ldrsb — one constructor (Ldrsb) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrsbParams — parameters of the ldrsb form.
type LdrsbParams struct {
	LsParams
}

func NewLdrsbParams(p LsParams) LdrsbParams {
	return LdrsbParams{LsParams: p}
}

func (p LdrsbParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldrsb(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrsbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldrsb — an arbitrary ldrsb.
func Ldrsb(rnd *rand.Rand) ohsnap.Arbitrary[LdrsbParams] {
	rtX := true
	base := newLsGen(rnd, &rtX, ScaledOff(rnd, 0))
	return ldrsbArb{base: base}
}

type ldrsbArb struct {
	base lsGen
}

func (a ldrsbArb) Generate() iter.Seq[LdrsbParams] {
	return arbStream(func() LdrsbParams {
		return NewLdrsbParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldrsbArb) Shrink(p LdrsbParams) iter.Seq[LdrsbParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]LdrsbParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdrsbParams(s))
	}

	return slices.Values(out)
}
