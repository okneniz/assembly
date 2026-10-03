package arm64

// Generator for msr — one generator, one type, one constructor (Msr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// MsrParams — parameters of msr sysreg, rt.
type MsrParams struct {
	Sysreg string
	Rt     arm64.Reg
}

func NewMsrParams(sysreg string, rt arm64.Reg) MsrParams {
	return MsrParams{
		Sysreg: sysreg,
		Rt:     rt,
	}
}

func (p MsrParams) Instr() arm64.Instr {
	in, err := arm64.New().Msr(p.Sysreg, p.Rt)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p MsrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// msrGen — generator for msr: rt is an x-register, the sysreg name is
// uniform over the arch table.
type msrGen struct {
	rnd     *rand.Rand
	sysregs ohsnap.Arbitrary[string]
}

// Msr — an arbitrary msr.
func Msr(rnd *rand.Rand) ohsnap.Arbitrary[MsrParams] {
	return newMsrGen(rnd)
}

func newMsrGen(rnd *rand.Rand) msrGen {
	return msrGen{
		rnd:     rnd,
		sysregs: Sysreg(rnd),
	}
}

func (g msrGen) Generate() iter.Seq[MsrParams] {
	return arb.Stream(func() MsrParams {
		return NewMsrParams(ohsnap.First(g.sysregs.Generate()), genReg(g.rnd, true, false, true))
	})
}

func (g msrGen) Shrink(p MsrParams) iter.Seq[MsrParams] {
	shrinks := slices.Collect(g.sysregs.Shrink(p.Sysreg))
	out := make([]MsrParams, 0, len(shrinks))
	for _, name := range shrinks {
		out = append(out, NewMsrParams(name, p.Rt))
	}

	for _, r := range regShrunk(p.Rt) {
		out = append(out, NewMsrParams(p.Sysreg, r))
	}

	return slices.Values(out)
}
