package arm64

// Generator for rev — one generator, one type, one constructor (Rev).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// RevParams — parameters of rev rd, rn (the 64-bit form only).
type RevParams struct {
	Rd, Rn arm64.Reg
}

func NewRevParams(rd arm64.Reg, rn arm64.Reg) RevParams {
	return RevParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p RevParams) Instr() arm64.Instr {
	in, err := arm64.New().Rev(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p RevParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// revGen — generator for rev: x-registers, occasionally xzr.
type revGen struct {
	rnd *rand.Rand
}

func newRevGen(rnd *rand.Rand) revGen {
	return revGen{rnd: rnd}
}

// Rev — an arbitrary rev.
func Rev(rnd *rand.Rand) ohsnap.Arbitrary[RevParams] {
	return newRevGen(rnd)
}

func (g revGen) Generate() iter.Seq[RevParams] {
	return arb.Stream(func() RevParams {
		return NewRevParams(genReg(g.rnd, true, false, true), genReg(g.rnd, true, false, true))
	})
}

func (g revGen) Shrink(p RevParams) iter.Seq[RevParams] {
	var out []RevParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewRevParams(r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewRevParams(p.Rd, r))
	}

	return slices.Values(out)
}
