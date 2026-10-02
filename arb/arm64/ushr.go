package arm64

// Generator for ushr — one constructor (Ushr) over the shared
// element-bound shift core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UshrParams — parameters of ushr.
type UshrParams struct {
	VShiftParams
}

func NewUshrParams(p VShiftParams) UshrParams {
	return UshrParams{VShiftParams: p}
}

func (p UshrParams) Instr() arm64.Instr {
	in, err := arm64.New().Ushr(p.Rd, p.Rn, p.Arr, p.Shift)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UshrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ushr — an arbitrary ushr.
func Ushr(rnd *rand.Rand) ohsnap.Arbitrary[UshrParams] {
	base := newVShiftGen1(rnd, arrFull())
	return ushrArb{base: base}
}

type ushrArb struct {
	base vShiftGen
}

func (a ushrArb) Generate() iter.Seq[UshrParams] {
	return arbStream(func() UshrParams {
		return NewUshrParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ushrArb) Shrink(p UshrParams) iter.Seq[UshrParams] {
	var out []UshrParams
	for _, s := range slices.Collect(a.base.Shrink(p.VShiftParams)) {
		out = append(out, NewUshrParams(s))
	}

	return slices.Values(out)
}
