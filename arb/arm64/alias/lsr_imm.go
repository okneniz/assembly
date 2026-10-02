package alias

// Generator for the lsr immediate alias — one generator, one type, one
// text form family: lsr rd, rn, #sh (the UBFM encoding; the register
// form is the LsrReg family of arb/arm64). The 64-bit form shifts
// 1..63, the 32-bit one 1..31 (the zero amount assembles as the
// canonical lsr #0 spelling of the mov degenerate — outside the
// family).

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// LsrImmParams — parameters of the lsr immediate alias.
type LsrImmParams struct {
	Rd, Rn arm64.Reg
	Sh     uint32
}

func NewLsrImmParams(rd arm64.Reg, rn arm64.Reg, sh uint32) LsrImmParams {
	return LsrImmParams{
		Rd: rd,
		Rn: rn,
		Sh: sh,
	}
}

func (p LsrImmParams) String() string {
	return fmt.Sprintf("lsr %s, %s, #%d", p.Rd, p.Rn, p.Sh)
}

func (p LsrImmParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// lsrImmGen — generator for lsr #imm: registers of the same width, the
// amount 1..regsize-1.
type lsrImmGen struct {
	rnd *rand.Rand
}

func newLsrImmGen(rnd *rand.Rand) lsrImmGen {
	return lsrImmGen{rnd: rnd}
}

// LsrImm — an arbitrary lsr #imm.
func LsrImm(rnd *rand.Rand) ohsnap.Arbitrary[LsrImmParams] {
	return newLsrImmGen(rnd)
}

func (g lsrImmGen) Generate() iter.Seq[LsrImmParams] {
	return stream(func() LsrImmParams {
		is64 := g.rnd.IntN(2) == 1
		sh := uint32(1 + g.rnd.IntN(62))
		if !is64 {
			sh = uint32(1 + g.rnd.IntN(30))
		}

		return NewLsrImmParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			sh,
		)
	})
}

func (g lsrImmGen) Shrink(p LsrImmParams) iter.Seq[LsrImmParams] {
	var out []LsrImmParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewLsrImmParams(r, p.Rn, p.Sh))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewLsrImmParams(p.Rd, r, p.Sh))
	}

	for _, v := range uhalved(p.Sh) {
		if v == 0 {
			continue // the zero amount is the mov degenerate
		}

		out = append(out, NewLsrImmParams(p.Rd, p.Rn, v))
	}

	return slices.Values(out)
}
