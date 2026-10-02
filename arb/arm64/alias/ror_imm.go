package alias

// Generator for the ror immediate alias — one generator, one type, one
// text form family: ror rd, rn, #imm (the EXTR encoding with rn == rm;
// the register form is the RorReg family of arb/arm64). The 64-bit
// form rotates 0..63, the 32-bit one 0..31.

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// RorImmParams — parameters of the ror immediate alias.
type RorImmParams struct {
	Rd, Rn arm64.Reg
	Sh     uint32
}

func NewRorImmParams(rd arm64.Reg, rn arm64.Reg, sh uint32) RorImmParams {
	return RorImmParams{
		Rd: rd,
		Rn: rn,
		Sh: sh,
	}
}

func (p RorImmParams) String() string {
	return fmt.Sprintf("ror %s, %s, #%d", p.Rd, p.Rn, p.Sh)
}

func (p RorImmParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// rorImmGen — generator for ror #imm: registers of the same width, the
// amount 0..regsize-1.
type rorImmGen struct {
	rnd *rand.Rand
}

func newRorImmGen(rnd *rand.Rand) rorImmGen {
	return rorImmGen{rnd: rnd}
}

// RorImm — an arbitrary ror #imm.
func RorImm(rnd *rand.Rand) ohsnap.Arbitrary[RorImmParams] {
	return newRorImmGen(rnd)
}

func (g rorImmGen) Generate() iter.Seq[RorImmParams] {
	return stream(func() RorImmParams {
		is64 := g.rnd.IntN(2) == 1
		sh := uint32(g.rnd.IntN(64))
		if !is64 {
			sh = uint32(g.rnd.IntN(32))
		}

		return NewRorImmParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			sh,
		)
	})
}

func (g rorImmGen) Shrink(p RorImmParams) iter.Seq[RorImmParams] {
	var out []RorImmParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewRorImmParams(r, p.Rn, p.Sh))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewRorImmParams(p.Rd, r, p.Sh))
	}

	for _, v := range uhalved(p.Sh) {
		out = append(out, NewRorImmParams(p.Rd, p.Rn, v))
	}

	return slices.Values(out)
}
