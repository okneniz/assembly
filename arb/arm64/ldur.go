package arm64

// Generator for ldur — one constructor (Ldur) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdurParams — parameters of the ldur form.
type LdurParams struct {
	LsParams
}

func NewLdurParams(p LsParams) LdurParams {
	return LdurParams{LsParams: p}
}

func (p LdurParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldur(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdurParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldur — an arbitrary ldur.
func Ldur(rnd *rand.Rand) ohsnap.Arbitrary[LdurParams] {
	base := newLsGen(rnd, nil, UnscaledOff(rnd))
	return ldurArb{base: base}
}

type ldurArb struct {
	base lsGen
}

func (a ldurArb) Generate() iter.Seq[LdurParams] {
	return arbStream(func() LdurParams {
		return NewLdurParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldurArb) Shrink(p LdurParams) iter.Seq[LdurParams] {
	var out []LdurParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdurParams(s))
	}

	return slices.Values(out)
}
