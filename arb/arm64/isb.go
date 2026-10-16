package arm64

// Generator for isb — one generator, one type, one constructor (Isb).

import (
	"iter"
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// IsbParams — parameters of isb (none: the bare spelling is Sy).
type IsbParams struct{}

func NewIsbParams() IsbParams {
	return IsbParams{}
}

func (p IsbParams) Instr() arm64.Instr {
	in, err := arm64.New().Isb()
	if err != nil {
		return nil // unreachable: isb has no operands to refuse
	}

	return in
}

func (p IsbParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// isbGen — generator for isb: a constant family.
type isbGen struct {
	rnd *rand.Rand
}

// Isb — an arbitrary isb.
func Isb(rnd *rand.Rand) ohsnap.Arbitrary[IsbParams] {
	return newIsbGen(rnd)
}

func newIsbGen(rnd *rand.Rand) isbGen {
	return isbGen{rnd: rnd}
}

func (g isbGen) Generate() IsbParams {
	return NewIsbParams()
}

func (g isbGen) Shrink(p IsbParams) iter.Seq[IsbParams] {
	return shrink.NoShrink[IsbParams]()(p)
}
