package arm64

// Generator for dmb — one generator, one type, one constructor (Dmb).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// DmbParams — parameters of dmb #domain.
type DmbParams struct {
	Domain arm64.BarrierDomain
}

func NewDmbParams(domain arm64.BarrierDomain) DmbParams {
	return DmbParams{Domain: domain}
}

func (p DmbParams) Instr() arm64.Instr {
	in, err := arm64.New().Dmb(p.Domain)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p DmbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// dmbGen — generator for dmb: the domain over the eight arch spellings.
type dmbGen struct {
	domains ohsnap.Arbitrary[arm64.BarrierDomain]
}

// Dmb — an arbitrary dmb.
func Dmb(rnd *rand.Rand) ohsnap.Arbitrary[DmbParams] {
	return newDmbGen(rnd)
}

func newDmbGen(rnd *rand.Rand) dmbGen {
	return dmbGen{domains: BarrierDomain(rnd)}
}

func (g dmbGen) Generate() iter.Seq[DmbParams] {
	return arb.Stream(func() DmbParams {
		return NewDmbParams(ohsnap.First(g.domains.Generate()))
	})
}

func (g dmbGen) Shrink(p DmbParams) iter.Seq[DmbParams] {
	shrinks := slices.Collect(g.domains.Shrink(p.Domain))
	out := make([]DmbParams, 0, len(shrinks))
	for _, d := range shrinks {
		out = append(out, NewDmbParams(d))
	}

	return slices.Values(out)
}
