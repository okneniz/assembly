package arm64

// Generator for fcvt — the width conversion: rd and rn of OPPOSITE kinds
// (d←s or s←d).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// FcvtParams — parameters of fcvt.
type FcvtParams struct {
	Rd, Rn arm64.FReg
}

func NewFcvtParams(rd arm64.FReg, rn arm64.FReg) FcvtParams {
	return FcvtParams{
		Rd: rd,
		Rn: rn,
	}
}

func (p FcvtParams) Instr() arm64.Instr {
	in, err := arm64.New().Fcvt(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p FcvtParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// fcvtGen — generator for fcvt: the kinds are opposite by construction.
type fcvtGen struct {
	rnd *rand.Rand
}

// Fcvt — an arbitrary fcvt.
func Fcvt(rnd *rand.Rand) ohsnap.Arbitrary[FcvtParams] {
	return newFcvtGen(rnd)
}

func newFcvtGen(rnd *rand.Rand) fcvtGen {
	return fcvtGen{rnd: rnd}
}

func (g fcvtGen) Generate() iter.Seq[FcvtParams] {
	return arbStream(func() FcvtParams {
		is64 := g.rnd.IntN(2) == 1
		return NewFcvtParams(genFp(g.rnd, is64), genFp(g.rnd, !is64))
	})
}

func (g fcvtGen) Shrink(p FcvtParams) iter.Seq[FcvtParams] {
	fps := fpShrunk(p.Rd)
	out := make([]FcvtParams, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewFcvtParams(r, p.Rn))
	}

	for _, r := range fpShrunk(p.Rn) {
		out = append(out, NewFcvtParams(p.Rd, r))
	}

	return slices.Values(out)
}
