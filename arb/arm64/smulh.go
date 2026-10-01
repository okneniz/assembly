package arm64

// Generator for smulh — one generator, one type, one constructor (Smulh).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SmulhParams — parameters of smulh rd, rn, rm (the 64-bit form only).
type SmulhParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewSmulhParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) SmulhParams {
	return SmulhParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p SmulhParams) Instr() arm64.Instr {
	in, err := arm64.New().Smulh(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SmulhParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// smulhGen — generator for smulh: x-registers, occasionally xzr.
type smulhGen struct {
	rnd *rand.Rand
}

func newSmulhGen(rnd *rand.Rand) smulhGen {
	return smulhGen{rnd: rnd}
}

// Smulh — an arbitrary smulh.
func Smulh(rnd *rand.Rand) ohsnap.Arbitrary[SmulhParams] {
	return newSmulhGen(rnd)
}

func (g smulhGen) Generate() iter.Seq[SmulhParams] {
	return arb.Stream(func() SmulhParams {
		return NewSmulhParams(
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
		)
	})
}

func (g smulhGen) Shrink(p SmulhParams) iter.Seq[SmulhParams] {
	var out []SmulhParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewSmulhParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewSmulhParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewSmulhParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
