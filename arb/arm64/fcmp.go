package arm64

// Generator for fcmp — one constructor (Fcmp) over the shared
// two-register FP core (one s or d kind).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/arbitrary"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FcmpParams — parameters of fcmp.
type FcmpParams struct {
	Rn, Rm arm64.FReg
}

func NewFcmpParams(rn arm64.FReg, rm arm64.FReg) FcmpParams {
	return FcmpParams{
		Rn: rn,
		Rm: rm,
	}
}

func (p FcmpParams) Instr() arm64.Instr {
	in, err := arm64.New().Fcmp(p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p FcmpParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// fcmpGen — generator for fcmp.
type fcmpGen struct {
	rnd  *rand.Rand
	coin ohsnap.Arbitrary[bool]
}

// Fcmp — an arbitrary fcmp.
func Fcmp(rnd *rand.Rand) ohsnap.Arbitrary[FcmpParams] {
	return newFcmpGen(rnd)
}

func newFcmpGen(rnd *rand.Rand) fcmpGen {
	return fcmpGen{
		rnd:  rnd,
		coin: arbitrary.Bool(rnd),
	}
}

func (g fcmpGen) Generate() FcmpParams {
	is64 := g.coin.Generate()
	return NewFcmpParams(genFp(g.rnd, is64), genFp(g.rnd, is64))
}

func (g fcmpGen) Shrink(p FcmpParams) iter.Seq[FcmpParams] {
	fps := fpShrunk(p.Rn)
	out := make([]FcmpParams, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewFcmpParams(r, p.Rm))
	}

	for _, r := range fpShrunk(p.Rm) {
		out = append(out, NewFcmpParams(p.Rn, r))
	}

	return slices.Values(out)
}
