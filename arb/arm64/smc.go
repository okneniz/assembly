package arm64

// Generator for smc — one generator, one type, one constructor (Smc).

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

// SmcParams — parameters of smc #imm16.
type SmcParams struct {
	Imm arm64.Imm16
}

func NewSmcParams(imm arm64.Imm16) SmcParams {
	return SmcParams{Imm: imm}
}

func (p SmcParams) Instr() arm64.Instr {
	return arm64.New().Smc(p.Imm)
}
func (p SmcParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// smcGen — generator for smc: immediate uniform in 0..0xffff.
type smcGen struct {
	rnd *rand.Rand
}

func newSmcGen(rnd *rand.Rand) smcGen {
	return smcGen{rnd: rnd}
}

// Smc — an arbitrary smc.
func Smc(rnd *rand.Rand) ohsnap.Arbitrary[SmcParams] {
	return newSmcGen(rnd)
}

func (g smcGen) Generate() iter.Seq[SmcParams] {
	return arb.Stream(func() SmcParams {
		return NewSmcParams(imm16(g.rnd.Int64N(0x10000)))
	})
}

func (g smcGen) Shrink(p SmcParams) iter.Seq[SmcParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[SmcParams]() // String() of our own type is unparseable — invariant
	}

	var out []SmcParams
	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm16(d)
		if err != nil {
			continue // unreachable: half of a valid imm16 is always in 0..65535
		}

		out = append(out, NewSmcParams(imm))
	}

	return slices.Values(out)
}
