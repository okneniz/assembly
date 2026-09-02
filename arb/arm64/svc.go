package arm64

// Generator for svc — one generator, one type, one constructor (Svc).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SvcParams — parameters of svc #imm16.
type SvcParams struct {
	Imm arm64.Imm16
}

func NewSvcParams(imm arm64.Imm16) SvcParams {
	return SvcParams{Imm: imm}
}

func (p SvcParams) Instr() arm64.Instr {
	return arm64.New().Svc(p.Imm)
}
func (p SvcParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// svcGen — generator for svc: immediate uniform in 0..0xffff.
type svcGen struct {
	rnd *rand.Rand
}

func newSvcGen(rnd *rand.Rand) svcGen {
	return svcGen{rnd: rnd}
}

// Svc — an arbitrary svc.
func Svc(rnd *rand.Rand) ohsnap.Arbitrary[SvcParams] {
	return newSvcGen(rnd)
}

func (g svcGen) Generate() iter.Seq[SvcParams] {
	return arb.Stream(func() SvcParams {
		return NewSvcParams(imm16(g.rnd.Int64N(0x10000)))
	})
}

func (g svcGen) Shrink(p SvcParams) iter.Seq[SvcParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[SvcParams]() // String() of our own type is unparseable — invariant
	}

	var out []SvcParams
	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm16(d)
		if err != nil {
			continue // unreachable: half of a valid imm16 is always in 0..65535
		}

		out = append(out, NewSvcParams(imm))
	}

	return slices.Values(out)
}
