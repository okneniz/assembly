package arm64

// Generator for add (extended register) — one constructor (AddExt) over
// the shared ext core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddExtParams — parameters of add rd, rn, rm, ext #imm3.
type AddExtParams struct {
	ExtParams
}

func NewAddExtParams(p ExtParams) AddExtParams {
	return AddExtParams{ExtParams: p}
}

func (p AddExtParams) Instr() arm64.Instr {
	in, err := arm64.New().AddExt(p.Rd, p.Rn, p.Rm, p.Ext, p.Imm3)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddExtParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// AddExt — an arbitrary add (extended register).
func AddExt(rnd *rand.Rand) ohsnap.Arbitrary[AddExtParams] {
	base := ext(rnd)
	return addExtArb{base: base}
}

type addExtArb struct {
	base extGen
}

func (a addExtArb) Generate() iter.Seq[AddExtParams] {
	return arbStream(func() AddExtParams {
		return NewAddExtParams(ohsnap.First(a.base.Generate()))
	})
}

func (a addExtArb) Shrink(p AddExtParams) iter.Seq[AddExtParams] {
	var out []AddExtParams
	for _, s := range slices.Collect(a.base.Shrink(p.ExtParams)) {
		out = append(out, NewAddExtParams(s))
	}

	return slices.Values(out)
}
