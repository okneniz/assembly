package arm64

// Generator for ldrsw — one constructor (Ldrsw) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrswParams — parameters of the ldrsw form.
type LdrswParams struct {
	LsParams
}

func NewLdrswParams(p LsParams) LdrswParams {
	return LdrswParams{LsParams: p}
}

func (p LdrswParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldrsw(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrswParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldrsw — an arbitrary ldrsw.
func Ldrsw(rnd *rand.Rand) ohsnap.Arbitrary[LdrswParams] {
	rtX := true
	base := newLsGen(rnd, &rtX, ScaledOff(rnd, 2))
	return ldrswArb{base: base}
}

type ldrswArb struct {
	base lsGen
}

func (a ldrswArb) Generate() iter.Seq[LdrswParams] {
	return arbStream(func() LdrswParams {
		return NewLdrswParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldrswArb) Shrink(p LdrswParams) iter.Seq[LdrswParams] {
	shrinks := slices.Collect(a.base.Shrink(p.LsParams))
	out := make([]LdrswParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdrswParams(s))
	}

	return slices.Values(out)
}
