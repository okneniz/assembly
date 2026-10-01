package arm64

// Generator for nop — one generator, one type, one constructor (Nop).

import (
	"iter"
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// NopParams — parameters of nop (none).
type NopParams struct{}

func NewNopParams() NopParams {
	return NopParams{}
}

func (p NopParams) Instr() arm64.Instr {
	return arm64.New().Nop()
}
func (p NopParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// nopGen — generator for nop: a constant family (compositions and the
// differential need it among the sources).
type nopGen struct {
	rnd *rand.Rand
}

func newNopGen(rnd *rand.Rand) nopGen {
	return nopGen{rnd: rnd}
}

// Nop — an arbitrary nop.
func Nop(rnd *rand.Rand) ohsnap.Arbitrary[NopParams] {
	return newNopGen(rnd)
}

func (g nopGen) Generate() iter.Seq[NopParams] {
	return arb.Stream(func() NopParams {
		return NewNopParams()
	})
}

func (g nopGen) Shrink(p NopParams) iter.Seq[NopParams] {
	return ohsnap.Empty[NopParams]()
}
