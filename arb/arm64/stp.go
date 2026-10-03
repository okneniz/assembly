package arm64

// Generator for stp — one constructor (Stp) over the shared pair core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StpParams — parameters of the stp form.
type StpParams struct {
	PairParams
}

func NewStpParams(p PairParams) StpParams {
	return StpParams{PairParams: p}
}

func (p StpParams) Instr() arm64.Instr {
	in, err := arm64.New().Stp(p.Rt, p.Rt2, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StpParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Stp — an arbitrary stp.
func Stp(rnd *rand.Rand) ohsnap.Arbitrary[StpParams] {
	base := pair(rnd)
	return stpArb{base: base}
}

type stpArb struct {
	base pairGen
}

func (a stpArb) Generate() iter.Seq[StpParams] {
	return arbStream(func() StpParams {
		return NewStpParams(ohsnap.First(a.base.Generate()))
	})
}

func (a stpArb) Shrink(p StpParams) iter.Seq[StpParams] {
	shrinks := slices.Collect(a.base.Shrink(p.PairParams))
	out := make([]StpParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewStpParams(s))
	}

	return slices.Values(out)
}
