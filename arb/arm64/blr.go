package arm64

// Generator for blr — one generator, one type, one constructor (Blr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BlrParams — parameters of blr rn.
type BlrParams struct {
	Rn arm64.Reg
}

func NewBlrParams(rn arm64.Reg) BlrParams {
	return BlrParams{Rn: rn}
}

func (p BlrParams) Instr() arm64.Instr {
	in, err := arm64.New().Blr(p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p BlrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// blrGen — generator for blr: rn is an x-register, occasionally xzr.
type blrGen struct {
	rnd *rand.Rand
}

// Blr — an arbitrary blr.
func Blr(rnd *rand.Rand) ohsnap.Arbitrary[BlrParams] {
	return newBlrGen(rnd)
}

func newBlrGen(rnd *rand.Rand) blrGen {
	return blrGen{rnd: rnd}
}

func (g blrGen) Generate() iter.Seq[BlrParams] {
	return arb.Stream(func() BlrParams {
		return NewBlrParams(genReg(g.rnd, true, false, true))
	})
}

func (g blrGen) Shrink(p BlrParams) iter.Seq[BlrParams] {
	rn := regShrunk(p.Rn)
	out := make([]BlrParams, 0, len(rn))
	for _, r := range rn {
		out = append(out, NewBlrParams(r))
	}

	return slices.Values(out)
}
