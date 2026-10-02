package alias

// Generator for the lsl immediate alias — one generator, one type, one
// text form family: lsl rd, rn, #sh (the UBFM encoding; the register
// form is the LslReg family of arb/arm64). The 64-bit form shifts
// 1..63, the 32-bit one 1..31 (the zero amount is the lsr #0/mov
// degenerate of the legacy path — outside the family).

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// LslImmParams — parameters of the lsl immediate alias.
type LslImmParams struct {
	Rd, Rn arm64.Reg
	Sh     uint32
}

func NewLslImmParams(rd arm64.Reg, rn arm64.Reg, sh uint32) LslImmParams {
	return LslImmParams{
		Rd: rd,
		Rn: rn,
		Sh: sh,
	}
}

func (p LslImmParams) String() string {
	return fmt.Sprintf("lsl %s, %s, #%d", p.Rd, p.Rn, p.Sh)
}

func (p LslImmParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// lslImmGen — generator for lsl #imm: registers of the same width, the
// amount 1..regsize-1.
type lslImmGen struct {
	rnd *rand.Rand
}

func newLslImmGen(rnd *rand.Rand) lslImmGen {
	return lslImmGen{rnd: rnd}
}

// LslImm — an arbitrary lsl #imm.
func LslImm(rnd *rand.Rand) ohsnap.Arbitrary[LslImmParams] {
	return newLslImmGen(rnd)
}

func (g lslImmGen) Generate() iter.Seq[LslImmParams] {
	return stream(func() LslImmParams {
		is64 := g.rnd.IntN(2) == 1
		sh := uint32(1 + g.rnd.IntN(62))
		if !is64 {
			sh = uint32(1 + g.rnd.IntN(30))
		}

		return NewLslImmParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			sh,
		)
	})
}

func (g lslImmGen) Shrink(p LslImmParams) iter.Seq[LslImmParams] {
	var out []LslImmParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewLslImmParams(r, p.Rn, p.Sh))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewLslImmParams(p.Rd, r, p.Sh))
	}

	for _, v := range uhalved(p.Sh) {
		if v == 0 {
			continue // the zero amount is the lsr #0 degenerate
		}

		out = append(out, NewLslImmParams(p.Rd, p.Rn, v))
	}

	return slices.Values(out)
}
