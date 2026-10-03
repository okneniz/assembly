package alias

// Generator for the neg alias — one generator, one type, one text form
// family: neg rd, rm[, shift #amt].

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// NegParams — parameters of the neg alias.
type NegParams struct {
	Rd, Rm arm64.Reg
	Sh     string // "", lsl/lsr/asr/ror
	Amt    int64
}

func NewNegParams(rd arm64.Reg, rm arm64.Reg, sh string, amt int64) NegParams {
	return NegParams{
		Rd:  rd,
		Rm:  rm,
		Sh:  sh,
		Amt: amt,
	}
}

func (p NegParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p NegParams) String() string {
	return "neg " + p.Rd.String() + ", " + p.Rm.String() + shiftSuffix(p.Sh, p.Amt)
}

// negGen — generator for neg: same-width registers, the shift amount
// below the register width.
type negGen struct {
	rnd *rand.Rand
}

// Neg — an arbitrary neg.
func Neg(rnd *rand.Rand) ohsnap.Arbitrary[NegParams] {
	return newNegGen(rnd)
}

func newNegGen(rnd *rand.Rand) negGen {
	return negGen{rnd: rnd}
}

func (g negGen) Generate() iter.Seq[NegParams] {
	return stream(func() NegParams {
		is64 := g.rnd.IntN(2) == 1
		width := int64(32)
		if is64 {
			width = 64
		}

		sh := ""
		amt := int64(0)
		if g.rnd.IntN(2) == 1 {
			sh = addSubShiftKinds()[g.rnd.IntN(3)]
			amt = g.rnd.Int64N(width-1) + 1 // #0 is the no-shift canonical spelling
		}

		return NewNegParams(
			a64.GenReg(
				g.rnd,
				is64,
				false,
				false,
			), // zr would read back as cmp (a decoder-canon gap)
			a64.GenReg(g.rnd, is64, false, true),
			sh,
			amt,
		)
	})
}

func (g negGen) Shrink(p NegParams) iter.Seq[NegParams] {
	var out []NegParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewNegParams(r, p.Rm, p.Sh, p.Amt))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewNegParams(p.Rd, r, p.Sh, p.Amt))
	}

	for _, v := range halved(p.Amt) {
		if v == 0 {
			continue // the no-shift form is the candidate below
		}

		out = append(out, NewNegParams(p.Rd, p.Rm, p.Sh, v))
	}

	if p.Sh != "" {
		out = append(out, NewNegParams(p.Rd, p.Rm, "", 0))
	}

	return slices.Values(out)
}
