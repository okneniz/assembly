package arm64

// Generator for rbit — one generator, one type, one constructor (Rbit).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/arbitrary"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// RbitParams — parameters of rbit rd, rn.
type RbitParams struct {
	Rd, Rn arm64.Reg
}

func NewRbitParams(rd arm64.Reg, rn arm64.Reg) RbitParams {
	return RbitParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p RbitParams) Instr() arm64.Instr {
	in, err := arm64.New().Rbit(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p RbitParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// rbitGen — generator for rbit: registers of the same width, 31st is zr.
type rbitGen struct {
	rnd  *rand.Rand
	coin ohsnap.Arbitrary[bool]
}

// Rbit — an arbitrary rbit.
func Rbit(rnd *rand.Rand) ohsnap.Arbitrary[RbitParams] {
	return newRbitGen(rnd)
}

func newRbitGen(rnd *rand.Rand) rbitGen {
	return rbitGen{
		rnd:  rnd,
		coin: arbitrary.Bool(rnd),
	}
}

func (g rbitGen) Generate() RbitParams {
	is64 := g.coin.Generate()
	return NewRbitParams(genReg(g.rnd, is64, false, true), genReg(g.rnd, is64, false, true))
}

func (g rbitGen) Shrink(p RbitParams) iter.Seq[RbitParams] {
	regs := regShrunk(p.Rd)
	out := make([]RbitParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewRbitParams(r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewRbitParams(p.Rd, r))
	}

	return slices.Values(out)
}
