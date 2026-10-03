package arm64

// Generator for br — one generator, one type, one constructor (Br).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BrParams — parameters of br rn.
type BrParams struct {
	Rn arm64.Reg
}

func NewBrParams(rn arm64.Reg) BrParams {
	return BrParams{Rn: rn}
}

func (p BrParams) Instr() arm64.Instr {
	in, err := arm64.New().Br(p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p BrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// brGen — generator for br: rn is an x-register, occasionally xzr.
type brGen struct {
	rnd *rand.Rand
}

// Br — an arbitrary br.
func Br(rnd *rand.Rand) ohsnap.Arbitrary[BrParams] {
	return newBrGen(rnd)
}

func newBrGen(rnd *rand.Rand) brGen {
	return brGen{rnd: rnd}
}

func (g brGen) Generate() iter.Seq[BrParams] {
	return arb.Stream(func() BrParams {
		return NewBrParams(genReg(g.rnd, true, false, true))
	})
}

func (g brGen) Shrink(p BrParams) iter.Seq[BrParams] {
	rn := regShrunk(p.Rn)
	out := make([]BrParams, 0, len(rn))
	for _, r := range rn {
		out = append(out, NewBrParams(r))
	}

	return slices.Values(out)
}
