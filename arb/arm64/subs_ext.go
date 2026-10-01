package arm64

// Generator for subs (extended register) — one constructor (SubsExt) over
// the shared ext core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SubsExtParams — parameters of subs rd, rn, rm, ext #imm3.
type SubsExtParams struct {
	ExtParams
}

func NewSubsExtParams(p ExtParams) SubsExtParams {
	return SubsExtParams{ExtParams: p}
}

func (p SubsExtParams) Instr() arm64.Instr {
	in, err := arm64.New().SubsExt(p.Rd, p.Rn, p.Rm, p.Ext, p.Imm3)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SubsExtParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// SubsExt — an arbitrary subs (extended register).
func SubsExt(rnd *rand.Rand) ohsnap.Arbitrary[SubsExtParams] {
	base := ext(rnd)
	return subsExtArb{base: base}
}

type subsExtArb struct {
	base extGen
}

func (a subsExtArb) Generate() iter.Seq[SubsExtParams] {
	return arbStream(func() SubsExtParams {
		return NewSubsExtParams(ohsnap.First(a.base.Generate()))
	})
}

func (a subsExtArb) Shrink(p SubsExtParams) iter.Seq[SubsExtParams] {
	var out []SubsExtParams
	for _, s := range slices.Collect(a.base.Shrink(p.ExtParams)) {
		out = append(out, NewSubsExtParams(s))
	}

	return slices.Values(out)
}
