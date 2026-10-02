package alias

// Generator for the asr immediate alias — one generator, one type, one
// text form family: asr rd, rn, #sh (the SBFM encoding; the register
// form is the AsrReg family of arb/arm64). The 64-bit form shifts
// 1..63, the 32-bit one 1..31 (the zero amount is the mov degenerate —
// outside the family).

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// AsrImmParams — parameters of the asr immediate alias.
type AsrImmParams struct {
	Rd, Rn arm64.Reg
	Sh     uint32
}

func NewAsrImmParams(rd arm64.Reg, rn arm64.Reg, sh uint32) AsrImmParams {
	return AsrImmParams{
		Rd: rd,
		Rn: rn,
		Sh: sh,
	}
}

func (p AsrImmParams) String() string {
	return fmt.Sprintf("asr %s, %s, #%d", p.Rd, p.Rn, p.Sh)
}

func (p AsrImmParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// asrImmGen — generator for asr #imm: registers of the same width, the
// amount 1..regsize-1.
type asrImmGen struct {
	rnd *rand.Rand
}

func newAsrImmGen(rnd *rand.Rand) asrImmGen {
	return asrImmGen{rnd: rnd}
}

// AsrImm — an arbitrary asr #imm.
func AsrImm(rnd *rand.Rand) ohsnap.Arbitrary[AsrImmParams] {
	return newAsrImmGen(rnd)
}

func (g asrImmGen) Generate() iter.Seq[AsrImmParams] {
	return stream(func() AsrImmParams {
		is64 := g.rnd.IntN(2) == 1
		sh := uint32(1 + g.rnd.IntN(62))
		if !is64 {
			sh = uint32(1 + g.rnd.IntN(30))
		}

		return NewAsrImmParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			sh,
		)
	})
}

func (g asrImmGen) Shrink(p AsrImmParams) iter.Seq[AsrImmParams] {
	var out []AsrImmParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewAsrImmParams(r, p.Rn, p.Sh))
	}

	for _, r := range a64.RegShrunk(p.Rn) {
		out = append(out, NewAsrImmParams(p.Rd, r, p.Sh))
	}

	for _, v := range uhalved(p.Sh) {
		if v == 0 {
			continue // the zero amount is the mov degenerate
		}

		out = append(out, NewAsrImmParams(p.Rd, p.Rn, v))
	}

	return slices.Values(out)
}
