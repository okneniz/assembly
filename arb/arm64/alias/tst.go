package alias

// Generator for the tst alias — one generator, one type, one text form
// family: tst rn, rm[, shift #amt] (the register form; the logical
// immediate form waits for the structural bitmask generator).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// TstParams — parameters of the tst alias (register form).
type TstParams struct {
	Rn, Rm arm64.Reg
	Sh     string // "", lsl/lsr/asr/ror
	Amt    int64
}

func NewTstParams(rn arm64.Reg, rm arm64.Reg, sh string, amt int64) TstParams {
	return TstParams{
		Rn:  rn,
		Rm:  rm,
		Sh:  sh,
		Amt: amt,
	}
}

func (p TstParams) String() string {
	return "tst " + p.Rn.String() + ", " + p.Rm.String() + shiftSuffix(p.Sh, p.Amt)
}

func (p TstParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// tstGen — generator for tst: same-width registers, the shift amount
// below the register width.
type tstGen struct {
	rnd *rand.Rand
}

func newTstGen(rnd *rand.Rand) tstGen {
	return tstGen{rnd: rnd}
}

// Tst — an arbitrary tst.
func Tst(rnd *rand.Rand) ohsnap.Arbitrary[TstParams] {
	return newTstGen(rnd)
}

func (g tstGen) Generate() iter.Seq[TstParams] {
	return stream(func() TstParams {
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

		return NewTstParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			sh,
			amt,
		)
	})
}

func (g tstGen) Shrink(p TstParams) iter.Seq[TstParams] {
	var out []TstParams
	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewTstParams(r, p.Rm, p.Sh, p.Amt))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewTstParams(p.Rn, r, p.Sh, p.Amt))
	}

	for _, v := range halved(p.Amt) {
		if v == 0 {
			continue // the no-shift form is the candidate below
		}

		out = append(out, NewTstParams(p.Rn, p.Rm, p.Sh, v))
	}

	if p.Sh != "" {
		out = append(out, NewTstParams(p.Rn, p.Rm, "", 0))
	}

	return slices.Values(out)
}
