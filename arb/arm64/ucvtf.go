package arm64

// Generator for ucvtf — one constructor (Ucvtf) over the shared
// conversion core (rd is the FP register, rn the gpr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UcvtfParams — parameters of ucvtf.
type UcvtfParams struct {
	FGprParams
}

func NewUcvtfParams(p FGprParams) UcvtfParams {
	return UcvtfParams{FGprParams: p}
}

func (p UcvtfParams) Instr() arm64.Instr {
	in, err := arm64.New().Ucvtf(p.F, p.R)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UcvtfParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ucvtf — an arbitrary ucvtf.
func Ucvtf(rnd *rand.Rand) ohsnap.Arbitrary[UcvtfParams] {
	base := fgprAny(rnd)
	return ucvtfArb{base: base}
}

type ucvtfArb struct {
	base fgprGen
}

func (a ucvtfArb) Generate() iter.Seq[UcvtfParams] {
	return arbStream(func() UcvtfParams {
		return NewUcvtfParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ucvtfArb) Shrink(p UcvtfParams) iter.Seq[UcvtfParams] {
	shrinks := slices.Collect(a.base.Shrink(p.FGprParams))
	out := make([]UcvtfParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewUcvtfParams(s))
	}

	return slices.Values(out)
}
