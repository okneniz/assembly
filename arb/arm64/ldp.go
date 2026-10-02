package arm64

// Generator for ldp — one constructor (Ldp) over the shared pair core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdpParams — parameters of the ldp form.
type LdpParams struct {
	PairParams
}

func NewLdpParams(p PairParams) LdpParams {
	return LdpParams{PairParams: p}
}

func (p LdpParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldp(p.Rt, p.Rt2, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdpParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldp — an arbitrary ldp.
func Ldp(rnd *rand.Rand) ohsnap.Arbitrary[LdpParams] {
	base := pair(rnd)
	return ldpArb{base: base}
}

type ldpArb struct {
	base pairGen
}

func (a ldpArb) Generate() iter.Seq[LdpParams] {
	return arbStream(func() LdpParams {
		return NewLdpParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldpArb) Shrink(p LdpParams) iter.Seq[LdpParams] {
	var out []LdpParams
	for _, s := range slices.Collect(a.base.Shrink(p.PairParams)) {
		out = append(out, NewLdpParams(s))
	}

	return slices.Values(out)
}
