package arm64

// Generator for sshr — one constructor (Sshr) over the shared
// element-bound shift core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SshrParams — parameters of sshr.
type SshrParams struct {
	VShiftParams
}

func NewSshrParams(p VShiftParams) SshrParams {
	return SshrParams{VShiftParams: p}
}

func (p SshrParams) Instr() arm64.Instr {
	in, err := arm64.New().Sshr(p.Rd, p.Rn, p.Arr, p.Shift)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SshrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Sshr — an arbitrary sshr.
func Sshr(rnd *rand.Rand) ohsnap.Arbitrary[SshrParams] {
	base := newVShiftGen1(rnd, arrFull())
	return sshrArb{base: base}
}

type sshrArb struct {
	base vShiftGen
}

func (a sshrArb) Generate() iter.Seq[SshrParams] {
	return arbStream(func() SshrParams {
		return NewSshrParams(ohsnap.First(a.base.Generate()))
	})
}

func (a sshrArb) Shrink(p SshrParams) iter.Seq[SshrParams] {
	var out []SshrParams
	for _, s := range slices.Collect(a.base.Shrink(p.VShiftParams)) {
		out = append(out, NewSshrParams(s))
	}

	return slices.Values(out)
}
