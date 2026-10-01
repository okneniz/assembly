package arm64

// Generator for mrs — one generator, one type, one constructor (Mrs).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// MrsParams — parameters of mrs rd, sysreg.
type MrsParams struct {
	Rd     arm64.Reg
	Sysreg string
}

func NewMrsParams(rd arm64.Reg, sysreg string) MrsParams {
	return MrsParams{
		Rd:     rd,
		Sysreg: sysreg,
	}
}

func (p MrsParams) Instr() arm64.Instr {
	in, err := arm64.New().Mrs(p.Rd, p.Sysreg)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p MrsParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// mrsGen — generator for mrs: rd is an x-register, the sysreg name is
// uniform over the arch table.
type mrsGen struct {
	rnd     *rand.Rand
	sysregs ohsnap.Arbitrary[string]
}

func newMrsGen(rnd *rand.Rand) mrsGen {
	return mrsGen{
		rnd:     rnd,
		sysregs: Sysreg(rnd),
	}
}

// Mrs — an arbitrary mrs.
func Mrs(rnd *rand.Rand) ohsnap.Arbitrary[MrsParams] {
	return newMrsGen(rnd)
}

func (g mrsGen) Generate() iter.Seq[MrsParams] {
	return arb.Stream(func() MrsParams {
		return NewMrsParams(genReg(g.rnd, true, false, true), ohsnap.First(g.sysregs.Generate()))
	})
}

func (g mrsGen) Shrink(p MrsParams) iter.Seq[MrsParams] {
	var out []MrsParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewMrsParams(r, p.Sysreg))
	}

	for _, name := range slices.Collect(g.sysregs.Shrink(p.Sysreg)) {
		out = append(out, NewMrsParams(p.Rd, name))
	}

	return slices.Values(out)
}
