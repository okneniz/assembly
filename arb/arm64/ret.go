package arm64

// Generator for ret — one generator, one type, one constructor (Ret).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// RetParams — parameters of ret rn.
type RetParams struct {
	Rn arm64.Reg
}

func NewRetParams(rn arm64.Reg) RetParams {
	return RetParams{Rn: rn}
}

func (p RetParams) Instr() arm64.Instr {
	in, err := arm64.New().Ret(p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p RetParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// retGen — generator for ret: rn is an x-register, occasionally xzr.
type retGen struct {
	rnd *rand.Rand
}

func newRetGen(rnd *rand.Rand) retGen {
	return retGen{rnd: rnd}
}

// Ret — an arbitrary ret.
func Ret(rnd *rand.Rand) ohsnap.Arbitrary[RetParams] {
	return newRetGen(rnd)
}

func (g retGen) Generate() iter.Seq[RetParams] {
	return arb.Stream(func() RetParams {
		return NewRetParams(genReg(g.rnd, true, false, true))
	})
}

func (g retGen) Shrink(p RetParams) iter.Seq[RetParams] {
	rn := regShrunk(p.Rn)
	out := make([]RetParams, 0, len(rn))
	for _, r := range rn {
		out = append(out, NewRetParams(r))
	}

	return slices.Values(out)
}
