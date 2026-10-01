package alias

// Generator for the mneg alias — one generator, one type, one text form
// family: mneg rd, rn, rm.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// MnegParams — parameters of the mneg alias.
type MnegParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewMnegParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) MnegParams {
	return MnegParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p MnegParams) String() string {
	return "mneg " + p.Rd.String() + ", " + p.Rn.String() + ", " + p.Rm.String()
}

func (p MnegParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// mnegGen — generator for mneg: same-width registers, 31st is zr.
type mnegGen struct {
	rnd *rand.Rand
}

func newMnegGen(rnd *rand.Rand) mnegGen {
	return mnegGen{rnd: rnd}
}

// Mneg — an arbitrary mneg.
func Mneg(rnd *rand.Rand) ohsnap.Arbitrary[MnegParams] {
	return newMnegGen(rnd)
}

func (g mnegGen) Generate() iter.Seq[MnegParams] {
	return stream(func() MnegParams {
		is64 := g.rnd.IntN(2) == 1
		return NewMnegParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
		)
	})
}

func (g mnegGen) Shrink(p MnegParams) iter.Seq[MnegParams] {
	var out []MnegParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewMnegParams(r, p.Rn, p.Rm))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewMnegParams(p.Rd, r, p.Rm))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewMnegParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
