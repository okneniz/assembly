package alias

// Generator for the mvn alias — one generator, one type, one text form
// family: mvn rd, rm[, shift #amt].

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// MvnParams — parameters of the mvn alias.
type MvnParams struct {
	Rd, Rm arm64.Reg
	Sh     string // "", lsl/lsr/asr/ror
	Amt    int64
}

func NewMvnParams(rd arm64.Reg, rm arm64.Reg, sh string, amt int64) MvnParams {
	return MvnParams{
		Rd:  rd,
		Rm:  rm,
		Sh:  sh,
		Amt: amt,
	}
}

func (p MvnParams) String() string {
	return "mvn " + p.Rd.String() + ", " + p.Rm.String() + shiftSuffix(p.Sh, p.Amt)
}

func (p MvnParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// mvnGen — generator for mvn: same-width registers, the shift amount
// below the register width.
type mvnGen struct {
	rnd *rand.Rand
}

func newMvnGen(rnd *rand.Rand) mvnGen {
	return mvnGen{rnd: rnd}
}

// Mvn — an arbitrary mvn.
func Mvn(rnd *rand.Rand) ohsnap.Arbitrary[MvnParams] {
	return newMvnGen(rnd)
}

func (g mvnGen) Generate() iter.Seq[MvnParams] {
	return stream(func() MvnParams {
		is64 := g.rnd.IntN(2) == 1
		width := int64(32)
		if is64 {
			width = 64
		}

		sh := ""
		amt := int64(0)
		if g.rnd.IntN(2) == 1 {
			sh = shiftKinds()[g.rnd.IntN(4)]
			amt = g.rnd.Int64N(width-1) + 1 // #0 is the no-shift canonical spelling
		}

		return NewMvnParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			sh,
			amt,
		)
	})
}

func (g mvnGen) Shrink(p MvnParams) iter.Seq[MvnParams] {
	var out []MvnParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewMvnParams(r, p.Rm, p.Sh, p.Amt))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewMvnParams(p.Rd, r, p.Sh, p.Amt))
	}

	for _, v := range halved(p.Amt) {
		if v == 0 {
			continue // the no-shift form is the candidate below
		}

		out = append(out, NewMvnParams(p.Rd, p.Rm, p.Sh, v))
	}

	if p.Sh != "" {
		out = append(out, NewMvnParams(p.Rd, p.Rm, "", 0))
	}

	return slices.Values(out)
}
