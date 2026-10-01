package arm64

// Generator for b.cond — one generator, one type, one constructor (Bcond).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BcondParams — parameters of b.cond (the condition and the pc-relative
// byte offset).
type BcondParams struct {
	Cond string
	Off  int64
}

func NewBcondParams(cond string, off int64) BcondParams {
	return BcondParams{
		Cond: cond,
		Off:  off,
	}
}

func (p BcondParams) Instr() arm64.Instr {
	in, err := arm64.New().Bcond(p.Cond, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BcondParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// bcondGen — generator for b.cond: a condition of the arch table and the
// offset uniform in the ±1MB imm19 range.
type bcondGen struct {
	rnd *rand.Rand
	off ohsnap.Arbitrary[int64]
}

func newBcondGen(rnd *rand.Rand) bcondGen {
	return bcondGen{
		rnd: rnd,
		off: BrOff(rnd, 1<<20),
	}
}

// Bcond — an arbitrary b.cond.
func Bcond(rnd *rand.Rand) ohsnap.Arbitrary[BcondParams] {
	return newBcondGen(rnd)
}

func (g bcondGen) Generate() iter.Seq[BcondParams] {
	return arb.Stream(func() BcondParams {
		return NewBcondParams(ohsnap.First(Cond(g.rnd).Generate()), ohsnap.First(g.off.Generate()))
	})
}

func (g bcondGen) Shrink(p BcondParams) iter.Seq[BcondParams] {
	var out []BcondParams
	for _, c := range slices.Collect(Cond(g.rnd).Shrink(p.Cond)) {
		out = append(out, NewBcondParams(c, p.Off))
	}

	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewBcondParams(p.Cond, v))
	}

	return slices.Values(out)
}
