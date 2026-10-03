package arm64

// Generator for ldpsw — one constructor (Ldpsw) over the shared pair core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdpswParams — parameters of the ldpsw form.
type LdpswParams struct {
	PairParams
}

func NewLdpswParams(p PairParams) LdpswParams {
	return LdpswParams{PairParams: p}
}

func (p LdpswParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldpsw(p.Rt, p.Rt2, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdpswParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldpsw — an arbitrary ldpsw.
func Ldpsw(rnd *rand.Rand) ohsnap.Arbitrary[LdpswParams] {
	base := pairX(rnd)
	return ldpswArb{base: base}
}

type ldpswArb struct {
	base pairGen
}

func (a ldpswArb) Generate() iter.Seq[LdpswParams] {
	return arbStream(func() LdpswParams {
		return NewLdpswParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldpswArb) Shrink(p LdpswParams) iter.Seq[LdpswParams] {
	shrinks := slices.Collect(a.base.Shrink(p.PairParams))
	out := make([]LdpswParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewLdpswParams(s))
	}

	return slices.Values(out)
}
