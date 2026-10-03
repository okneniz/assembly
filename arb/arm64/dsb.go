package arm64

// Generator for dsb — one generator, one type, one constructor (Dsb).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// DsbParams — parameters of dsb #domain.
type DsbParams struct {
	Domain arm64.BarrierDomain
}

func NewDsbParams(domain arm64.BarrierDomain) DsbParams {
	return DsbParams{Domain: domain}
}

func (p DsbParams) Instr() arm64.Instr {
	in, err := arm64.New().Dsb(p.Domain)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p DsbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// dsbGen — generator for dsb: the domain over the eight arch spellings.
type dsbGen struct {
	domains ohsnap.Arbitrary[arm64.BarrierDomain]
}

// Dsb — an arbitrary dsb.
func Dsb(rnd *rand.Rand) ohsnap.Arbitrary[DsbParams] {
	return newDsbGen(rnd)
}

func newDsbGen(rnd *rand.Rand) dsbGen {
	return dsbGen{domains: BarrierDomain(rnd)}
}

func (g dsbGen) Generate() iter.Seq[DsbParams] {
	return arb.Stream(func() DsbParams {
		return NewDsbParams(ohsnap.First(g.domains.Generate()))
	})
}

func (g dsbGen) Shrink(p DsbParams) iter.Seq[DsbParams] {
	shrinks := slices.Collect(g.domains.Shrink(p.Domain))
	out := make([]DsbParams, 0, len(shrinks))
	for _, d := range shrinks {
		out = append(out, NewDsbParams(d))
	}

	return slices.Values(out)
}
