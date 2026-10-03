package arm64

// Generator for sub (extended register) — one constructor (SubExt) over
// the shared ext core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SubExtParams — parameters of sub rd, rn, rm, ext #imm3.
type SubExtParams struct {
	ExtParams
}

func NewSubExtParams(p ExtParams) SubExtParams {
	return SubExtParams{ExtParams: p}
}

func (p SubExtParams) Instr() arm64.Instr {
	in, err := arm64.New().SubExt(p.Rd, p.Rn, p.Rm, p.Ext, p.Imm3)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SubExtParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// SubExt — an arbitrary sub (extended register).
func SubExt(rnd *rand.Rand) ohsnap.Arbitrary[SubExtParams] {
	base := ext(rnd)
	return subExtArb{base: base}
}

type subExtArb struct {
	base extGen
}

func (a subExtArb) Generate() iter.Seq[SubExtParams] {
	return arbStream(func() SubExtParams {
		return NewSubExtParams(ohsnap.First(a.base.Generate()))
	})
}

func (a subExtArb) Shrink(p SubExtParams) iter.Seq[SubExtParams] {
	shrinks := slices.Collect(a.base.Shrink(p.ExtParams))
	out := make([]SubExtParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewSubExtParams(s))
	}

	return slices.Values(out)
}
