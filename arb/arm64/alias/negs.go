package alias

// Generator for the negs alias — one generator, one type, one text form
// family: negs rd, rm[, shift #amt].

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// NegsParams — parameters of the negs alias.
type NegsParams struct {
	Rd, Rm arm64.Reg
	Sh     string // "", lsl/lsr/asr/ror
	Amt    int64
}

func NewNegsParams(rd arm64.Reg, rm arm64.Reg, sh string, amt int64) NegsParams {
	return NegsParams{
		Rd:  rd,
		Rm:  rm,
		Sh:  sh,
		Amt: amt,
	}
}

func (p NegsParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p NegsParams) String() string {
	return "negs " + p.Rd.String() + ", " + p.Rm.String() + shiftSuffix(p.Sh, p.Amt)
}

// negsGen — generator for negs: same-width registers, the shift amount
// below the register width.
type negsGen struct {
	rnd *rand.Rand
}

// Negs — an arbitrary negs.
func Negs(rnd *rand.Rand) ohsnap.Arbitrary[NegsParams] {
	return newNegsGen(rnd)
}

func newNegsGen(rnd *rand.Rand) negsGen {
	return negsGen{rnd: rnd}
}

func (g negsGen) Generate() iter.Seq[NegsParams] {
	return stream(func() NegsParams {
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

		return NewNegsParams(
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

func (g negsGen) Shrink(p NegsParams) iter.Seq[NegsParams] {
	var out []NegsParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewNegsParams(r, p.Rm, p.Sh, p.Amt))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewNegsParams(p.Rd, r, p.Sh, p.Amt))
	}

	for _, v := range halved(p.Amt) {
		if v == 0 {
			continue // the no-shift form is the candidate below
		}

		out = append(out, NewNegsParams(p.Rd, p.Rm, p.Sh, v))
	}

	if p.Sh != "" {
		out = append(out, NewNegsParams(p.Rd, p.Rm, "", 0))
	}

	return slices.Values(out)
}
